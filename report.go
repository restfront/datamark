package datamark

type ReportRequest struct {
	ReportID *string `json:"report_id,omitempty"`

	Group  *string       `json:"group,omitempty"`
	Labels []string      `json:"labels,omitempty"`
	Params []ReportParam `json:"params,omitempty"`
}

type ReportParam struct {
	Code  int    `json:"code"`
	Value string `json:"value"`
}

type ReportResponse struct {
	UID       string `json:"report_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Type      int    `json:"type"`
	Group     string `json:"group"`

	Status ReportStatus `json:"status"`
	Result ReportResult `json:"result"`
}

type ReportStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ReportResult struct {
	Labels    ReportLabel `json:"labels"`
	TypeID    int         `json:"type_id"`
	ErrorCode int         `json:"error"`
	Message   string      `json:"message"`
}

type ReportLabel struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Details []string `json:"details"`
}
