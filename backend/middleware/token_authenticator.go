package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"
)

var secretKey = []byte("just a random secret string")

func TokenAuthenticator() gin.HandlerFunc {
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
			// log err
			context.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid authentication token. Error: %v.", err))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user, present := claims["userID"]
			if !present {
				context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
				return
			}

			context.Set("userID", user)
		} else {
			context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
			return
		}

		context.Next()
	}
}
