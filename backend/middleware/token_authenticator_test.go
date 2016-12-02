package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gin-gonic/gin.v1"
)

type mockSecretProvider struct{}

func (msp *mockSecretProvider) GetSecretKey() []byte {
	return []byte("secret")
}

func TestAllGood(t *testing.T) {
	secretProvider := &mockSecretProvider{}
	called := false

	correctJWT := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6NDEwMjQ0NDgwMH0.7yRo-iMpG-tp01AAKpvtQAm1ZbX1o6L4n1h5Wws0snw"

	router := gin.New()
	router.Use(TokenAuthenticator(secretProvider))
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
	secretProvider := &mockSecretProvider{}
	called := false

	unsupportedSigningMethodToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6NDEwMjQ0NDgwMH0.im_HO5T6Oy-n7kvwQvJIFrakpqr1IqAngFtQ4FUjiaY"

	router := gin.New()
	router.Use(TokenAuthenticator(secretProvider))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+unsupportedSigningMethodToken)

	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWrongSignature(t *testing.T) {
	secretProvider := &mockSecretProvider{}
	called := false

	unsupportedSigningMethodToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6NDEwMjQ0NDgwMH0.DqWD84gFXT2reQB064MourQ5RT4lhXreEhEEcibWSZQ"

	router := gin.New()
	router.Use(TokenAuthenticator(secretProvider))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+unsupportedSigningMethodToken)

	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestExpiredToken(t *testing.T) {
	secretProvider := &mockSecretProvider{}
	called := false

	expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6OTQ2Njg0ODAwfQ.qCQVAYbj2G0zba0tjiq4bBfjqRqyjKtEh_YD-KAexC4"

	router := gin.New()
	router.Use(TokenAuthenticator(secretProvider))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestNoExpiryDateToken(t *testing.T) {
	secretProvider := &mockSecretProvider{}
	called := false

	expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjF9.ZcqhDVtPl_Qm71xVxhGuVJVxDBA7ifm1IXOhJTe-FPc"

	router := gin.New()
	router.Use(TokenAuthenticator(secretProvider))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestNoUserIDToken(t *testing.T) {
	secretProvider := &mockSecretProvider{}
	called := false

	expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQxMDI0NDQ4MDB9.GKBB3kfb1XwV7KKCX5RaLBn48U453La4IlmgtKSU6So"

	router := gin.New()
	router.Use(TokenAuthenticator(secretProvider))
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.String(200, "%d", c.MustGet("userID").(int64))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
