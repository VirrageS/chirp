package token

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"

	"github.com/VirrageS/chirp/backend/config"
	serviceErrors "github.com/VirrageS/chirp/backend/model/errors"
)

type tokenManager struct {
	secretKey                  []byte
	authTokenValidityPeriod    time.Duration
	refreshTokenValidityPeriod time.Duration
}

func NewManager(config config.TokenConfigProvider) Manager {
	return &tokenManager{
		secretKey:                  config.GetSecretKey(),
		authTokenValidityPeriod:    config.GetAuthTokenValidityPeriod(),
		refreshTokenValidityPeriod: config.GetRefreshTokenValidityPeriod(),
	}
}

func (m *tokenManager) ValidateToken(tokenString string, request *http.Request) (int64, error) {
	// set up a parser that doesn't validate expiration time
	parser := jwt.Parser{}
	parser.SkipClaimsValidation = true

	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
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

		// check if requester IP is correct
		if err := m.verifyIP(claims, request); err != nil {
			return 0, err
		}

		// chcek if userAgent is correct
		if err := m.verifyUserAgent(claims, request); err != nil {
			return 0, err
		}

		return int64(userID), nil
	}

	return 0, errors.New("Malformed authentication token.")
}

func (m *tokenManager) CreateAuthToken(userID int64, request *http.Request) (string, error) {
	return m.createToken(userID, request, m.authTokenValidityPeriod)
}

func (m *tokenManager) CreateRefreshToken(userID int64, request *http.Request) (string, error) {
	return m.createToken(userID, request, m.refreshTokenValidityPeriod)
}

func (m *tokenManager) createToken(userID int64, request *http.Request, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	clientIP, err := m.getIPFromRequest(request)
	if err != nil {
		log.WithError(err).Fatal("Failed to sign token, error getting client IP from request.")
		return "", serviceErrors.UnexpectedError
	}
	userAgent := request.UserAgent()
	if userAgent == "" {
		return "", serviceErrors.NoUserAgentHeaderError
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":           userID,
		"allowedIP":        clientIP,
		"allowedUserAgent": userAgent,
		"exp":              expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		log.WithError(err).Fatal("Failed to sign token, error signing the token.")
		return "", serviceErrors.UnexpectedError
	}

	return tokenString, nil
}

func (m *tokenManager) verifyIP(claims jwt.MapClaims, request *http.Request) error {
	requestIP, err := m.getIPFromRequest(request)
	if err != nil {
		return fmt.Errorf("Malformed request: %s.", err)
	}

	claimsExpectedIP, isSetAllowedIP := claims["allowedIP"]
	expectedIP, ok := claimsExpectedIP.(string)
	if !ok || !isSetAllowedIP {
		return errors.New("Token does not contain required data: allowedIP.")
	}

	if requestIP != expectedIP {
		return errors.New("Token is not allowed to be used from this IP.")
	}

	return nil
}

func (m *tokenManager) verifyUserAgent(claims jwt.MapClaims, request *http.Request) error {
	requestUserAgent := request.UserAgent()
	if requestUserAgent == "" {
		return errors.New("Malformed request: no User-Agent header.")
	}

	claimsExpectedUserAgent, isSetUserAgent := claims["allowedUserAgent"]
	expectedUserAgent, ok := claimsExpectedUserAgent.(string)
	if !ok || !isSetUserAgent {
		return errors.New("Token does not contain required data: allowedUserAgent.")
	}

	if requestUserAgent != expectedUserAgent {
		return errors.New("Token is not allowed to be used from this User-Agent.")
	}

	return nil
}

func (m *tokenManager) getIPFromRequest(request *http.Request) (string, error) {
	// We expect the realclient IP to be in X-Real-Ip header
	if headerIP := request.Header.Get("X-Real-Ip"); headerIP != "" {
		return headerIP, nil
	}

	// If IP was not present in X-Real-Ip header, we assume that RemoteAddr is the real client IP
	remoteAddr := request.RemoteAddr
	if remoteAddr == "" {
		log.WithField("request", fmt.Sprintf("%+v", request)).
			Error("No RemoteAddr and X-Real-IP header in request in verifyIP.")

		return "", errors.New("no client IP was provided")
	}

	headerIP, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"request": fmt.Sprintf("%+v", request),
			"address": remoteAddr,
		}).Error("Error splitting RemoteAddr header value into host:port.")

		return "", errors.New("invalid RemoteAddr format, expected host:port format")
	}

	return headerIP, nil
}
