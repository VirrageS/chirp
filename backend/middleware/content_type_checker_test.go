package middleware

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/gin-gonic/gin.v1"
)

var _ = Describe("ContentTypeChecker", func() {
	var (
		router *gin.Engine
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		router.Use(ContentTypeChecker())

		router.POST("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
	})

	It("should allow to make normal request when content type header is set", func() {
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/test", nil)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, request)

		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("should not allow to make request when content type header is set but has wrong value", func() {
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/test", nil)
		request.Header.Set("Content-Type", "application/xml")

		router.ServeHTTP(w, request)

		Expect(w.Code).To(Equal(http.StatusUnsupportedMediaType))
	})

	It("should not allow to make request when content type header is not set", func() {
		w := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/test", nil)

		router.ServeHTTP(w, request)

		Expect(w.Code).To(Equal(http.StatusUnsupportedMediaType))
	})
})
