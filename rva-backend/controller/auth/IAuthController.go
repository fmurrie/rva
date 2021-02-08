package auth

type IAuthController interface {
	GenerateToken(information interface{}) string
	ValidateToken(token string) bool
}
