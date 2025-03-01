package viewmodel

type AuthLoginVM struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AuthTokenVM struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshVM struct {
	RefreshToken string `json:"refresh_token" `
}
