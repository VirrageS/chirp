package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/VirrageS/chirp/backend/token"
)

// TokenAuthenticator check if token is valid and sets context key and value
// appropiretly. Aborts when there was a problem in validating token.
func TokenAuthenticator(tokenManager token.Manager) gin.HandlerFunc {
	return func(context *gin.Context) {
		fullTokenString := context.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(fullTokenString, "Bearer ")

		userID, err := tokenManager.ValidateToken(tokenString, context.Request)
		if err != nil {
			context.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		context.Set("userID", userID)
		context.Next()
	}
}
