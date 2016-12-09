package token

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/VirrageS/chirp/backend/config"
	serviceErrors "github.com/VirrageS/chirp/backend/model/errors"
)

type TokenManager struct {
	secretKey []byte
}

func NewTokenManager(secretKey config.SecretKeyProvider) TokenManagerProvider {
	return &TokenManager{
		secretKey: secretKey.GetSecretKey(),
	}
}

func (manager *TokenManager) ValidateToken(tokenString string) (int64, error) {
	// set up a parser that doesn't validate expiration time
	parser := jwt.Parser{}
	parser.SkipClaimsValidation = true

	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return manager.secretKey, nil
	})

	if err != nil {
		log.WithError(err).WithField("token", tokenString).Error("Failed to parse the token.")
		return 0, errors.New("Invalid authentication token.")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimUserID, isSetID := claims["userID"]
		userID, ok := claimUserID.(float64)
		if !ok || !isSetID {
			return 0, errors.New("Token does not contain required data.")
		}

		// check if token contains expiry date
		if unexpired := claims.VerifyExpiresAt(time.Now().Unix(), true); !unexpired {
			return 0, errors.New("Token has expired.")
		}

		return int64(userID), nil
	}

	return 0, errors.New("Malformed authentication token.")
}

func (manager *TokenManager) CreateToken(userID int64, duration int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(duration) * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(manager.secretKey)
	if err != nil {
		log.WithError(err).Fatal("Failed to sign token.")
		return "", serviceErrors.UnexpectedError
	}

	return tokenString, nil
}
