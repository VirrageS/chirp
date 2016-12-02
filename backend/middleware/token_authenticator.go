package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/config"
	"time"
)

func TokenAuthenticator(configuration config.SecretKeyProvider) gin.HandlerFunc {
	secretKey := configuration.GetSecretKey()

	return func(context *gin.Context) {
		fullTokenString := context.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(fullTokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil {
			log.WithError(err).WithField("token", tokenString).Error("Failed to parse the token.")
			context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			claimUserID, isSetID := claims["userID"]
			userID, ok := claimUserID.(float64)
			if !ok || !isSetID {
				context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
				return
			}

			// check if token contains expiry date
			if unexpired := claims.VerifyExpiresAt(time.Now().Unix(), true); !unexpired {
				context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
				return
			}

			context.Set("userID", int64(userID))
		} else {
			context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
			return
		}

		context.Next()
	}
}
