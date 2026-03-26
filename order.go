package datamark

type OrderRequest struct {
	ID *string `json:"id,omitempty"`

	Group  *string          `json:"group,omitempty"`
	Orders map[string]Order `json:"orders,omitempty"`

	// используется в /v2/orders/statusList
	LabelType *int `json:"label_type,omitempty"`
	Page      *int `json:"page,omitempty"`
}

type OrderResponse struct {
	Order Order `json:"order"`
}

type OrderGroupResponse map[string]Order

type OrdersListResponse struct {
	Orders []Order `json:"orders_list"`
}

type Order struct {
	ID          int    `json:"id"` // NOTE может быть и стринг
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Count       int    `json:"count"`
	Comment     string `json:"comment"`
	ParentOrder int    `json:"parent_order"`
	GTIN        string `json:"gtin"`
	CompletedAt string `json:"completed_at"`

	Status OrderStatus `json:"status"`
	Type   OrderType   `json:"type"`
	File   OrderFile   `json:"file"`
	User   User        `json:"user"`

	// используется в /v2/orders/addGroupOrders
	LabelType     int    `json:"label_type,omitempty"`
	TypographyID  int    `json:"typography_id,omitempty"`
	TypographyDoc string `json:"typography_doc,omitempty"` // Номер и дата заказа (договора с типографией)

	Labels []string `json:"labels"`
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
	Filename  string `json:"filename"`
	Downloads int    `json:"downloads"`
}
