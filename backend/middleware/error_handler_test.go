package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/square/go-jose.v1/json"
)

type errorResponse struct {
	Errors []string `json:"errors"`
}

type normalResponse struct {
	Response string `json:"response"`
}

func TestErrorOccurred(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.AbortWithError(400, errors.New("An error occured."))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	var resp errorResponse
	json.NewDecoder(w.Result().Body).Decode(&resp)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, errorResponse{[]string{"An error occured."}}, resp)
}

func TestNoErrorOccurred(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"response": "Response",
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	var resp normalResponse
	json.NewDecoder(w.Result().Body).Decode(&resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, normalResponse{"Response"}, resp)
}
