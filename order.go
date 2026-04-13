package datamark

type OrderRequest struct {
	ID *string `json:"id,omitempty"`

	// используется в /v3/orders/addGroupOrders
	Group  *string          `json:"group,omitempty"`
	Orders map[string]Order `json:"orders,omitempty"`

	// используется в /v2/orders/statusList
	LabelType *int `json:"label_type,omitempty"`
	Page      *int `json:"page,omitempty"`
}

type OrderResponse struct {
	Order Order `json:"detail"`
}

type OrderGroupResponse struct {
	Orders map[string]Order `json:"orders"`
}

type OrdersListResponse struct {
	Orders []Order `json:"orders_list"`
}

// Информация о заказе кодов маркировки
type Order struct {
	ID          int64  `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CompletedAt string `json:"completed_at"`
	Count       int    `json:"count"`
	Comment     string `json:"comment"`
	ParentOrder int    `json:"parent_order"`

	// GTIN string `json:"gtin"` // используется в GetOrdersStatuses v2 просто строка

	GTIN     OrderGTIN   `json:"gtin"`
	Group    OrderGroup  `json:"group"`
	Status   OrderStatus `json:"status"`
	Type     OrderType   `json:"type"`
	File     OrderFile   `json:"files"`
	User     User        `json:"user"`
	UserInfo AgentInfo   `json:"user_info"`

	// используется в /v3/orders/addGroupOrders
	LabelType     int    `json:"label_type"`
	TypographyID  int    `json:"typography_id,omitempty"`
	TypographyDoc string `json:"typography_doc,omitempty"` // Номер и дата заказа (договора с типографией)

	Labels []string `json:"labels"`
}

type OrderGTIN struct {
	Name    string `json:"name"`
	Articul string `json:"articul"`
	GTIN    string `json:"gtin"`
}

type OrderStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type OrderType struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Method       string `json:"method"`
	DigitsCount  int    `json:"digits_count"`
	IsBso        bool   `json:"is_bso"`
	LettersCount int    `json:"letters_count"`
	IsAutotake   bool   `json:"is_autotake"`
}

type OrderFile struct {
	Labels struct {
		Filename  string `json:"filename"`
		Downloads int    `json:"downloads"`
	} `json:"labels"`
}

type OrderGroup struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type LabelRequest struct {
	Filename string `json:"filename"`
}

type LabelResponse struct {
	Labels []string `json:"labels"`
}
