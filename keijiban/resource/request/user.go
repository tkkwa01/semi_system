package request

type UserCreate struct {
	Email        string `json:"email"`
	Introduction string `json:"introduction"`
}

type UserUpdate struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Introduction string `json:"introduction"`
}
