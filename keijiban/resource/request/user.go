package request

type UserCreate struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserUpdate struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Introduction string `json:"introduction"`
}

type UserLogin struct {
	Session  bool   `json:"session"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRefreshToken struct {
	Session      bool   `json:"session"`
	RefreshToken string `json:"refresh_token"`
}
