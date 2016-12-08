package token

type TokenManagerProvider interface {
	ValidateToken(tokenString string) (int64, error)
	CreateToken(userID int64, duration int) (string, error)
}
