package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/token"
)

type mockSecretProvider struct{}

func (msp *mockSecretProvider) GetSecretKey() []byte {
	return []byte("secret")
}

func mockTokenManagerProvider() token.TokenManagerProvider {
	return token.NewTokenManager(&mockSecretProvider{})
}

func TestAllGood(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tokenManager := mockTokenManagerProvider()
	called := false

	correctJWT := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6NDEwMjQ0NDgwMH0.7yRo-iMpG-tp01AAKpvtQAm1ZbX1o6L4n1h5Wws0snw"

	router := gin.New()
	router.Use(ErrorHandler())
	router.Use(TokenAuthenticator(tokenManager))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+correctJWT)

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "1", w.Body.String())
}

func TestUnsupportedSigningMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tokenManager := mockTokenManagerProvider()
	called := false

	unsupportedSigningMethodToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6NDEwMjQ0NDgwMH0.DqWD84gFXT2reQB064MourQ5RT4lhXreEhEEcibWSZQ"

	router := gin.New()
	router.Use(ErrorHandler())
	router.Use(TokenAuthenticator(tokenManager))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+unsupportedSigningMethodToken)

	router.ServeHTTP(w, req)

	var resp errorResponse
	json.NewDecoder(w.Body).Decode(&resp)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, errorResponse{[]string{"Invalid authentication token."}}, resp)
}

func TestWrongSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tokenManager := mockTokenManagerProvider()
	called := false

	wrongSignatureToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQxMDI0NDQ4MDB9.xfVyuL08DPhDgKSzIIeXnUwNjoSicw6MeMzSW3qVfM4"

	router := gin.New()
	router.Use(ErrorHandler())
	router.Use(TokenAuthenticator(tokenManager))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+wrongSignatureToken)

	router.ServeHTTP(w, req)

	var resp errorResponse
	json.NewDecoder(w.Body).Decode(&resp)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, errorResponse{[]string{"Invalid authentication token."}}, resp)
}

func TestExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tokenManager := mockTokenManagerProvider()
	called := false

	expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6OTQ2Njg0ODAwfQ.qCQVAYbj2G0zba0tjiq4bBfjqRqyjKtEh_YD-KAexC4"

	router := gin.New()
	router.Use(ErrorHandler())
	router.Use(TokenAuthenticator(tokenManager))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	router.ServeHTTP(w, req)

	var resp errorResponse
	json.NewDecoder(w.Body).Decode(&resp)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, errorResponse{[]string{"Token has expired."}}, resp)
}

// Test the case when token doesn't contain 'exp' field. This is incorrect and should be rejected.
func TestTokenWithNoExpiryDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tokenManager := mockTokenManagerProvider()
	called := false

	expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjF9.ZcqhDVtPl_Qm71xVxhGuVJVxDBA7ifm1IXOhJTe-FPc"

	router := gin.New()
	router.Use(ErrorHandler())
	router.Use(TokenAuthenticator(tokenManager))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	router.ServeHTTP(w, req)

	var resp errorResponse
	json.NewDecoder(w.Body).Decode(&resp)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, errorResponse{[]string{"Token has expired."}}, resp)
}

// Test the case when token doesn't contain 'userID' field. This is incorrect and should be rejected.
func TestTokenWithNoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tokenManager := mockTokenManagerProvider()
	called := false

	noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQxMDI0NDQ4MDB9.JtAP_exnIC2Dpdw7q_VCqI06vlzntvNsejr806cqBwA"

	router := gin.New()
	router.Use(ErrorHandler())
	router.Use(TokenAuthenticator(tokenManager))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+noUserIDToken)

	router.ServeHTTP(w, req)

	var resp errorResponse
	json.NewDecoder(w.Body).Decode(&resp)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, errorResponse{[]string{"Token does not contain required data."}}, resp)
}
