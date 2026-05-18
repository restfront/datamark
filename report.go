package datamark

type ReportRequest struct {
	UUIDs []string `json:"uuid_list,omitempty"`

	Group  *string     `json:"group,omitempty"`
	Labels []Label     `json:"labels,omitempty"`
	Params ReportParam `json:"params,omitempty"`
}

type Label struct {
	Label string `json:"label"`
}

type ReportParam struct {
	ManufactureDate          string `json:"manufacture_date,omitempty"`           // (25) Дата изготовления (производства) товара
	ItemName                 string `json:"item_name,omitempty"`                  // (94) Наименование товара
	MarkMethod               string `json:"mark_method,omitempty"`                // (95) Способ маркировки (полиграфическая защита)
	DeclarationRegNumber     string `json:"declaration_reg_number,omitempty"`     // (96) Регистрационный номер декларации на товары
	DeclarationReleaseDate   string `json:"declaration_release_date,omitempty"`   // (97) Дата выпуска товаров по декларации
	Country                  string `json:"country,omitempty"`                    // (100) Страна экспорта
	MarkTarget               string `json:"mark_target,omitempty"`                // (101) Цель маркировки
	ContractDate             string `json:"dogovor_date,omitempty"`               // (102) Основание: дата документа
	ContractNumber           string `json:"dogovor_nomer,omitempty"`              // (103) Основание: номер документа
	MarkReason               string `json:"mark_reason,omitempty"`                // (104) Причина нанесения СИ
	DocumentDate             string `json:"document_date,omitempty"`              // (109) Дата документа, подтверждающего приобретение товара
	DocumentNumber           string `json:"document_number,omitempty"`            // (110) Номер документа, подтверждающего приобретение товара
	TaxPayerNumber           string `json:"org_number,omitempty"`                 // (111) Номер налогоплательщика
	SellerName               string `json:"org_name,omitempty"`                   // (112) Наименование организации-продавца
	RemarkReason             string `json:"remark_reason,omitempty"`              // (124) Причина перемаркировки
	ExportCountryEAEU        string `json:"export_country_eaeu,omitempty"`        // (126) Страна экспорта (ЕАЭС)
	DeclarationOrdinalNumber string `json:"declaration_ordinal_number,omitempty"` // (134) Порядковый номер товара в декларации на товары
}

type ReportMarkResponse struct {
	ID   int64  `json:"report_id"`  // ID отчёта
	UUID string `json:"report_uid"` // UID, уникальный идентификатор отчёта о маркировке (report_type = 1)
}

type ReportResponseList struct {
	Reports []ReportResponse `json:"results"`
}

type ReportResponse struct {
	UUID      string `json:"report_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Type      int    `json:"type"`

	Status ReportStatus `json:"status"`
	Result ReportResult `json:"result"`

	NotFoundUUIDs []string `json:"not_found_uuids"`
}

type ReportStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ReportResult struct {
	Group  string      `json:"group"`
	GTIN   string      `json:"gtin"`
	Labels ReportLabel `json:"labels"`

	TypeID    int    `json:"type_id"`
	ErrorCode int    `json:"error"`
	Message   string `json:"message"`
}

type ReportLabel struct {
	Success      int      `json:"success"`
	Failed       int      `json:"failed"`
	Details      []string `json:"details"`
	FullQuantity int      `json:"all"`
}
