package token

import "net/http"

type TokenManagerProvider interface {
	ValidateToken(tokenString string, request *http.Request) (int64, error)
	CreateAuthToken(userID int64, request *http.Request) (string, error)
	CreateRefreshToken(userID int64, request *http.Request) (string, error)
}
