package datamark

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	timeFormat string = "2006-01-02 15:04:05"
	batchSize  int    = 100
)

func (client *Client) Auth(ctx context.Context) error {
	path := "/auth"

	endpoint, err := url.JoinPath(client.config.BaseURL, path)
	if err != nil {
		return ErrIncorrectURL
	}

	result := &AuthResponse{}
	errResponse := &ErrorResponse{}

	client.logger.Debugf("Запрос на авторизацию")

	response, err := client.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"username": client.config.Username,
			"password": client.config.Password,
		}).
		SetResult(result).
		SetError(errResponse).
		Post(endpoint)

	if err != nil {
		if isTimeout(err) {
			return ErrConnectionTimeout
		}
		return fmt.Errorf("ошибка при выполнении запроса авторизации: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("ошибка при выполнении запроса авторизации: %s", errResponse.ErrorDescription())
	}

	tokenExpTime, err := time.Parse(timeFormat, result.User.ExpiresIn)
	if err != nil {
		return fmt.Errorf("не удалось преобразовать время истечения токена: %w", err)
	}

	token := &accessToken{
		AccessToken: result.Token,
		ExpiresAt:   tokenExpTime.Unix(),
	}

	client.token = token

	return nil
}

func (client *Client) Logout(ctx context.Context) error {
	path := "/logout"

	client.mu.Lock()
	defer client.mu.Unlock()

	if client.token == nil || client.token.isExpired() {
		return nil
	}

	_, err := client.doRequest(ctx, http.MethodPost, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("ошибка при запросе выхода: %w", err)
	}

	client.token = nil

	return nil
}

// CreateItemByGtin добавление товара по GTIN
func (client *Client) CreateItemByGtin(ctx context.Context, request ItemRequest) (*ItemResponse, error) {
	path := "/items/addByGtin"

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на добавление товара по GTIN")

	result := &ItemResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе добавления товара по GTIN: %w", err)
	}

	return result, nil
}

// GetGTINStatuses проверка статусов регистрации GTIN.
// Ограничение: Не более 100 GTIN на 1 запрос
func (client *Client) GetGTINStatuses(ctx context.Context, request ItemRequest) (GTINStatusesResponse, error) {
	path := "/items/checkGtin"

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на проверку статусов регистрации GTIN")

	result := make(GTINStatusesResponse)
	for i := 0; i < len(request.GTINList); i += batchSize {
		end := i + batchSize
		// если длина списка меньше лимита, то конец списка будет остаток
		// или если изначально длина была меньше лимита, то конец списка будет длина основного списка
		if end > len(request.GTINList) {
			end = len(request.GTINList)
		}

		req := request.GTINList[i:end]

		_, err := client.doRequest(ctx, http.MethodPost, path, nil, req, &result)
		if err != nil {
			return result, fmt.Errorf("ошибка при запросе проверки статусов регистрации GTIN: %w", err)
		}
	}

	return result, nil
}

// FindItem поиск товара
// В старом коде этот метод ничего не делает с телом ответа, метод НЕ используется
func (client *Client) FindItem(ctx context.Context, request ItemRequest) (*ItemsResponse, error) {
	path := "/items/findItems"

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на поиск товара")

	result := &ItemsResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе поиска товара: %w", err)
	}

	// NOTE в этом методе в ответе поле gtin_check имеет тип bool.

	return result, nil
}

// GetOrderByID информация о заказе
func (client *Client) GetOrderByID(ctx context.Context, id int64) (*OrderResponse, error) {
	path := fmt.Sprintf("%s/%d", "/v3/orders/list", id)

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на получение информации о заказе")

	result := &OrderResponse{}
	_, err := client.doRequest(ctx, http.MethodGet, path, nil, nil, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе получения информации о заказе кодов маркировки: %w", err)
	}

	return result, nil
}

// CreateGroupOrders групповой заказ кодов маркировки
func (client *Client) CreateGroupOrders(ctx context.Context, request OrderRequest) (*OrderGroupResponse, error) {
	path := "/v3/orders/addGroupOrders"

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на создание группового заказа кодов маркировки")

	result := &OrderGroupResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе создания группового заказа: %w", err)
	}

	return result, nil
}

// GetOrdersStatuses Список статусов заказов
func (client *Client) GetOrdersStatuses(ctx context.Context, request OrderRequest) (*OrdersListResponse, error) {
	path := "/v2/orders/statusList"

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на получение списка статусов заказов")

	result := &OrdersListResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе получения списка статусов заказов: %w", err)
	}

	return result, nil
}

// GetLabels Список заказанных кодов
func (client *Client) GetLabels(ctx context.Context, request LabelRequest) (*LabelResponse, error) {
	path := "/v3/orders/downloads"

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на получение списка заказанных кодов")

	result := &LabelResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе получения списка заказанных кодов: %w", err)
	}

	return result, nil
}

// CreateMarksReport отчёт о маркировке
func (client *Client) CreateMarksReport(ctx context.Context, request ReportRequest) (*ReportMarkResponse, error) {
	path := "/v3/reports/addMark"

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на создание отчёта о маркировке")

	result := &ReportMarkResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе создания отчёта о маркировке: %w", err)
	}

	return result, nil
}

// GetReportStatus статус выполнения отчёта
func (client *Client) GetReportStatus(ctx context.Context, request ReportRequest) (*ReportResponseList, error) {
	path := "/v3/reports"

	if client.token == nil || client.token.isExpired() {
		err := client.refreshToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	client.logger.Debugf("Запрос на получение статусов выполнения отчёта")

	result := &ReportResponseList{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе получения статусов выполнения отчёта: %w", err)
	}

	return result, nil
}
