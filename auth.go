package datamark

type AuthRequest struct {
	Username     string
	Password     string
	IsRulesAgree bool
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"token_refresh"`
	User         User   `json:"user"`
}

type User struct {
	ID                 int    `json:"id"` // Код (идентификатор) пользователя
	Email              string `json:"email"`
	InfoID             int    `json:"info_id"` // ID субъекта (участника)
	Name               string `json:"name"`
	Lastname           string `json:"lastname"`
	NeedChangePassword bool   `json:"need_change_password"` // Признак необходимости смены пароля: False – не требуется
	Status             int    `json:"status"`               // 0 - заблокирован, 1 - активен
	ExpiresIn          string `json:"expires_in"`           // Срок действия токена

	Info AgentInfo `json:"info"`
	Role UserRole  `json:"role"`
}

// AgentInfo cведения о контрагенте (участнике)
type AgentInfo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	UNP             string `json:"unp"`
	GLN             string `json:"gln"`
	Country         string `json:"country"` // Страна регистрации контрагента
	Email           string `json:"email"`
	IsPartlyBlocked bool   `json:"is_partly_blocked"`
	IsRulesAgree    bool   `json:"is_rules_agree"` // Признак ознакомления и принятия условий публичного договора с Оператором на реализацию СИ и оказания услуг по предоставлению и учёту КМ в ГИС «Электронный знак»
	Status          int    `json:"status"`         // 0 - заблокирован, 1 - активен
}

type UserRole struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
