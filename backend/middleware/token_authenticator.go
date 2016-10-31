package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/config"
	"time"
)

var secretKey = config.GetSecretKey()

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
		// TODO: log err
		//fmt.Printf("Error parsing the token: = %v", err)
		context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimUserID, isSetID := claims["userID"]
		userID, ok := claimUserID.(float64)
		if !ok || !isSetID {
			//fmt.Printf("ok = %d, isSetID = %d\n", ok, isSetID)
			context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
			return
		}

		if unexpired := claims.VerifyExpiresAt(time.Now().Unix(), true); !unexpired {
			context.AbortWithError(http.StatusUnauthorized, errors.New("Authentication token has expired."))
			return
		}

		context.Set("userID", int64(userID))
	} else {
		//fmt.Printf("ok = %d, token.Valid = %d\n", ok, token.Valid)
		context.AbortWithError(http.StatusUnauthorized, errors.New("Invalid authentication token."))
		return
	}

	context.Next()
}
