package auth

type AuthType string

const (
	AuthTypeLocal AuthType = "local"
	AuthTypeOAuth AuthType = "oauth"
)



type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}



