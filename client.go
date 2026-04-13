package datamark

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	defaultRequestTimeout = time.Second * 15
	defaultMaxAttempts    = 3
	defaultDelaySec       = 1
)

type Client struct {
	config     *Config
	token      *accessToken
	httpClient *resty.Client
	logger     logger
	mu         sync.Mutex
}

type Config struct {
	BaseURL           string
	Username          string
	Password          string
	RequestTimeoutSec int
	RetryPolicy       struct {
		MaxAttempts int
		DelaySec    int
	}
	DebugMode bool
	DevMode   bool
}

type accessToken struct {
	AccessToken string
	ExpiresAt   int64
}

type logger interface {
	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
}

type noOpLogger struct{}

func (l *noOpLogger) Debugf(format string, v ...any) {}
func (l *noOpLogger) Infof(format string, v ...any)  {}
func (l *noOpLogger) Warnf(format string, v ...any)  {}
func (l *noOpLogger) Errorf(format string, v ...any) {}

func NewClient(config *Config, logger logger) *Client {
	if isNilInterface(logger) {
		logger = &noOpLogger{}
	}

	client := &Client{
		config: config,
		logger: logger,
		token:  nil,
		mu:     sync.Mutex{},
	}

	httpClient := resty.New()

	timeout := defaultRequestTimeout
	if config.RequestTimeoutSec > 0 {
		timeout = time.Duration(config.RequestTimeoutSec) * time.Second
	}

	maxAttemps := defaultMaxAttempts
	if config.RetryPolicy.MaxAttempts > 0 {
		maxAttemps = config.RetryPolicy.MaxAttempts
	}

	delaySec := defaultDelaySec
	if config.RetryPolicy.DelaySec > 0 {
		delaySec = config.RetryPolicy.DelaySec
	}

	httpClient.SetTimeout(timeout)
	httpClient.SetRetryCount(maxAttemps)
	httpClient.SetRetryWaitTime(time.Duration(delaySec) * time.Second)
	httpClient.SetHeader("User-Agent", "RestFront/datamark-client")

	if config.DebugMode {
		httpClient.SetDebug(true)
		httpClient.SetLogger(logger)
	}

	client.httpClient = httpClient

	return client
}

func isNilInterface(i interface{}) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)

	return v.Kind() == reflect.Ptr && v.IsNil()
}

func (c *Client) doRequest(
	ctx context.Context,
	method string,
	path string,
	queryParams url.Values,
	body any,
	result any,
) (*resty.Response, error) {
	// формирование URL
	endpoint, err := url.JoinPath(c.config.BaseURL, path)
	if err != nil {
		return nil, ErrIncorrectURL
	}

	errResponse := &ErrorResponse{}

	// инициализация запроса
	req := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Token", c.token.AccessToken).
		SetResult(result).
		SetError(errResponse)

	if len(queryParams) > 0 {
		req.SetQueryParamsFromValues(queryParams)
	}

	// заголовок и тело запроса
	if body != nil {
		req.SetHeader("Content-Type", "application/json").
			SetBody(body)
	}

	// выполнение запроса
	response := &resty.Response{}

	switch method {
	case http.MethodGet:
		response, err = req.Get(endpoint)
	case http.MethodPost:
		response, err = req.Post(endpoint)
	default:
		return nil, fmt.Errorf("%w: %s", ErrIncorrectRequestMethod, method)
	}

	// обработка ошибок
	if err != nil {
		if isTimeout(err) {
			return nil, ErrConnectionTimeout
		}
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %s", errResponse.ErrorDescription())
	}

	return response, nil
}

func retry(ctx context.Context, delay time.Duration, fn func(context.Context) (bool, error)) error {
	for {
		ok, err := fn(ctx)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", ErrConnectionTimeout, ctx.Err())
		}
	}
}

// refreshToken обновляет токен
func (client *Client) refreshToken(ctx context.Context) error {
	client.mu.Lock()
	defer client.mu.Unlock()

	return client.Auth(ctx)
}

// isExpired возвращает true если до окончания срока действия токена осталось меньше 5 минут
func (token *accessToken) isExpired() bool {
	return time.Until(time.Unix(token.ExpiresAt, 0)) < 5*time.Minute
}
