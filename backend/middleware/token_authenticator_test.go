package middleware

import (
	"net/http"
	"net/http/httptest"

	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = Describe("TokenAuthenticator", func() {
	var (
		router *gin.Engine
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)

		tokenManager := mockTokenManagerProvider()
		router = gin.New()
		router.Use(ErrorHandler())
		router.Use(TokenAuthenticator(tokenManager))
	})

	It("should allow to make normal response since jwt is okay", func() {
		correctJWT := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6NDEwMjQ0NDgwMH0.7yRo-iMpG-tp01AAKpvtQAm1ZbX1o6L4n1h5Wws0snw"

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+correctJWT)

		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(Equal("1"))
	})

	It("should return status unauthorized when token has unsupported signing method", func() {
		unsupportedSigningMethodToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6NDEwMjQ0NDgwMH0.DqWD84gFXT2reQB064MourQ5RT4lhXreEhEEcibWSZQ"

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+unsupportedSigningMethodToken)
		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Invalid authentication token."}}))
	})

	It("shoud return status unauthorized when token has wrong signature", func() {
		wrongSignatureToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQxMDI0NDQ4MDB9.xfVyuL08DPhDgKSzIIeXnUwNjoSicw6MeMzSW3qVfM4"

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+wrongSignatureToken)
		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Invalid authentication token."}}))
	})

	It("should return status unauthorized when token expired", func() {
		expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImV4cCI6OTQ2Njg0ODAwfQ.qCQVAYbj2G0zba0tjiq4bBfjqRqyjKtEh_YD-KAexC4"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token has expired."}}))
	})

	// Test the case when token does not contain 'exp' field. This is incorrect and should be rejected.
	It("should return status unauthorized when token does not contain `exp` field", func() {
		expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjF9.ZcqhDVtPl_Qm71xVxhGuVJVxDBA7ifm1IXOhJTe-FPc"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token has expired."}}))
	})

	// Test the case when token does not contain 'userID' field. This is incorrect and should be rejected.
	It("should return status unauthorized when token does not cotain `userID` field", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQxMDI0NDQ4MDB9.JtAP_exnIC2Dpdw7q_VCqI06vlzntvNsejr806cqBwA"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token does not contain required data."}}))
	})
})
