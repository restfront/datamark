package datamark

type ItemRequest struct {
	GTIN *string `json:"gtin,omitempty"`

	Search   *string `json:"search,omitempty"`    // Значение для поиска (GTIN либо Артикул/Модель)
	IsImport *bool   `json:"is_import,omitempty"` // Признак поиска в каталоге импортируемых товаров

	GTINList []string `json:"gtin_list,omitempty"`
}

type GTINStatusesResponse map[string]GTINStatus

type GTINStatus struct {
	Status    bool   `json:"gtin_check"` // Признак регистрации GTIN в каталоге: true – товар зарегистрирован; false – товар ожидает регистрации
	ErrorCode int    `json:"error"`
	Message   string `json:"message"`
}

type ItemsResponse struct {
	Items []ItemResponse `json:"items"`
}

type ItemResponse struct {
	Name        string `json:"name"`
	Image       string `json:"image"`   // URL-ссылка на изображение товара
	Articul     string `json:"articul"` // Артикул товара
	Code        string `json:"code"`    // Внутренний код (идентификатор) записи сведений о товаре
	Description string `json:"description"`
	GTIN        string `json:"gtin"`
	GTINCheck   any    `json:"gtin_check"` // NOTE Признак регистрации GTIN в каталоге: true – товар зарегистрирован; false – товар ожидает регистрации. Может быть bool и int
	Group       string `json:"group"`
	DocumentID  any    `json:"document_id"` // NOTE type?
	Tnved       string `json:"tnved"`       // Код ТН ВЭД ЕАЭС
	IsMyItem    bool   `json:"is_my"`       // Признак собственного товара: true – товар внесен текущим участником; false – товар внесен другим участником

	DebugMessage *ItemDebugMessage `json:"debug_message,omitempty"`
	Catalog      ItemCatalog       `json:"catalog"`
	Params       []ItemParam       `json:"params"`
}

type ItemDebugMessage struct {
	GTIN          string `json:"gtin"`
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

type ItemCatalog struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type ItemParam struct {
	Name         string `json:"name"`
	Units        any    `json:"units"`
	Code         int    `json:"code"`
	Value        string `json:"value"`
	DisplayValue string `json:"display_value"`
}
