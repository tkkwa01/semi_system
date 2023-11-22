package request

type UserCreate struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserUpdate struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
}

type UserLogin struct {
	Session  bool   `json:"session"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserRefreshToken struct {
	Session      bool   `json:"session"`
	RefreshToken string `json:"refresh_token"`
}
