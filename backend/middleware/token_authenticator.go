package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/config"
)

var secretKey = []byte(config.GetSecretKey())

func TokenAuthenticator(context *gin.Context) {
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
		context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimUserID, present := claims["userID"]
		userID, ok := claimUserID.(float64)

		if !present || !ok {
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
