package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/square/go-jose.v1/json"
)

type errorResponse struct {
	Errors []string `json:"errors"`
}

type normalResponse struct {
	Response string `json:"response"`
}

var _ = Describe("ErrorHandler", func() {
	var (
		router *gin.Engine
	)

	JustBeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		router.Use(ErrorHandler())
	})

	It("should propagate error", func() {
		router.GET("/test", func(c *gin.Context) {
			c.AbortWithError(http.StatusBadRequest, errors.New("An error occured."))
		})

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, request)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusBadRequest))
		Expect(response).To(Equal(errorResponse{[]string{"An error occured."}}))
	})

	It("should not fire error when everything is okay", func() {
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"response": "Response",
			})
		})

		w := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/test", nil)

		router.ServeHTTP(w, request)

		var response normalResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(response).To(Equal(normalResponse{"Response"}))
	})
})
