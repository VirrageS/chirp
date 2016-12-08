package middleware

import (
	"net/http"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/service"
)

func TokenAuthenticator(configuration config.SecretKeyProvider) gin.HandlerFunc {
	secretKey := configuration.GetSecretKey()

	return func(context *gin.Context) {
		fullTokenString := context.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(fullTokenString, "Bearer ")

		userID, err := service.ValidateToken(tokenString, secretKey)
		if err != nil {
			context.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		context.Set("userID", userID)
		context.Next()
	}
}
