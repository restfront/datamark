package datamark

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (client *Client) Auth(ctx context.Context, request AuthRequest) (*AuthResponse, error) {
	path := "/auth"

	if request.Username == "" {
		request.Username = client.config.DefaultUsername
	}

	if request.Password == "" {
		request.Password = client.config.DefaultPassword
	}

	endpoint, err := url.JoinPath(client.config.BaseURL, path)
	if err != nil {
		return nil, ErrIncorrectURL
	}

	result := &AuthResponse{}
	errResponse := &ErrorResponse{}

	client.logger.Debugf("Запрос на авторизацию")

	response, err := client.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"username": request.Username,
			"password": request.Password,
		}).
		SetResult(result).
		SetError(errResponse).
		Post(endpoint)

	if err != nil {
		if isTimeout(err) {
			return nil, ErrConnectionTimeout
		}
		return nil, fmt.Errorf("ошибка при выполнении запроса авторизации: %w", err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("ошибка при выполнении запроса авторизации: %s", errResponse.ErrorDescription())
	}

	client.token = result.Token

	return result, nil
}

func (client *Client) Logout(ctx context.Context) error {
	path := "/logout"

	_, err := client.doRequest(ctx, http.MethodPost, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("ошибка при запросе выхода: %w", err)
	}

	return nil
}

// CreateItemByGtin добавление товара по GTIN
func (client *Client) CreateItemByGtin(ctx context.Context, request ItemRequest) (*ItemResponse, error) {
	path := "/items/addByGtin"

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

	client.logger.Debugf("Запрос на проверку статусов регистрации GTIN")

	// TODO нужно ли учитывать ограничение в 100 GTIN, в старом коде отправляется список без проверки длины

	result := make(GTINStatusesResponse)
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе проверки статусов регистрации GTIN: %w", err)
	}

	return result, nil
}

// FindItem поиск товара
func (client *Client) FindItem(ctx context.Context, request ItemRequest) (*ItemsResponse, error) {
	path := "/items/findItems"

	client.logger.Debugf("Запрос на поиск товара")

	result := &ItemsResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе поиска товара: %w", err)
	}

	// NOTE в этом методе в ответе поле gtin_check имеет тип bool.
	// В старом коде этот метод ничего не делает с телом ответа, он вообще не используется

	return result, nil
}

// GetOrderByID информация о заказе
func (client *Client) GetOrderByID(ctx context.Context, request OrderRequest) (*OrderResponse, error) {
	path := "/v2/orders"

	client.logger.Debugf("Запрос на получение информации о заказе")

	result := &OrderResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе получения информации о заказе: %w", err)
	}

	// NOTE в этом методе в теле ответа поле ID имеет тип int

	return result, nil
}

// CreateGroupOrders групповой заказ кодов маркировки
func (client *Client) CreateGroupOrders(ctx context.Context, request OrderRequest) (OrderGroupResponse, error) {
	path := "/v2/orders/addGroupOrders"

	client.logger.Debugf("Запрос на создание группового заказа кодов маркировки")

	result := make(OrderGroupResponse)
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе создания группового заказа: %w", err)
	}

	// NOTE в этом методе в теле ответа поле ORDER.ID имеет тип string, либо сделать маленькую отдельную структуру с 2 полями

	return result, nil
}

// GetOrdersStatuses Список статусов заказов
func (client *Client) GetOrdersStatuses(ctx context.Context, request OrderRequest) (*OrdersListResponse, error) {
	path := "/v2/orders/statusList"

	client.logger.Debugf("Запрос на получение списка статусов заказов")

	result := &OrdersListResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе получения списка статусов заказов: %w", err)
	}

	// NOTE в этом методе в теле ответа поле ORDER.ID имеет тип int

	return result, nil
}

// GetFile скачивание файла
func (client *Client) GetFile(ctx context.Context, filename string) error {
	path := fmt.Sprintf("%s%s", "/downloads/", filename)

	client.logger.Debugf("Запрос на получение файла")

	_, err := client.doRequest(ctx, http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("ошибка при запросе получения файла: %w", err)
	}

	return nil
}

// CreateMarksReport отчёт о маркировке
func (client *Client) CreateMarksReport(ctx context.Context, request ReportRequest) (*ReportResponse, error) {
	path := "/v2/reports/addMark"

	client.logger.Debugf("Запрос на создание отчёта о маркировке")

	result := &ReportResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе создания отчёта о маркировке: %w", err)
	}

	return result, nil
}

// GetReportStatus статус выполнения отчёта
func (client *Client) GetReportStatus(ctx context.Context, request ReportRequest) (*ReportResponse, error) {
	path := "/v2/reports"

	client.logger.Debugf("Запрос на получение статуса выполнения отчёта")

	result := &ReportResponse{}
	_, err := client.doRequest(ctx, http.MethodPost, path, nil, request, result)
	if err != nil {
		return result, fmt.Errorf("ошибка при запросе получения статуса выполнения отчёта: %w", err)
	}

	return result, nil
}
