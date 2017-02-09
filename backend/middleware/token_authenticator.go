package middleware

import (
	"net/http"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/token"
)

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
