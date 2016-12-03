package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gin-gonic/gin.v1"
)

func TestSendJSONContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false

	router := gin.New()
	router.Use(ContentTypeChecker())
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestSendUnsupportedContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	called := false

	router := gin.New()
	router.Use(ContentTypeChecker())
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/xml")

	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, w.Code, http.StatusUnsupportedMediaType)
}

func TestNoContentType(t *testing.T) {
	called := false

	router := gin.New()
	router.Use(ContentTypeChecker())
	router.POST("/test", func(c *gin.Context) {
		called = true
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)

	router.ServeHTTP(w, req)

	assert.False(t, called)
	assert.Equal(t, w.Code, http.StatusUnsupportedMediaType)
}
