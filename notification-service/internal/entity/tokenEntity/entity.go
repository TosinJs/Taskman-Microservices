package tokenEntity

type IDTokenPayload struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
}
