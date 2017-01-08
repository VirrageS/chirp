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

func (msp *mockSecretProvider) GetAuthTokenValidityPeriod() int {
	return 1
}

func (msp *mockSecretProvider) GetRefreshTokenValidityPeriod() int {
	return 1
}

// From Go documentation of httptest.NewRequest():
// 192.0.2.0/24 is "TEST-NET" in RFC 5737 for use solely in
// documentation and example source code and should not be
// used publicly.
const (
	testIP = "192.0.2.1"
)

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

	It("should allow to make normal response when jwt token is okay", func() {
		correctJWT := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImV4cCI6NDEwMjQ0NDgwMH0.4svEQGu2aNhKryzBiAOTTaygvomLzPrA8yEe_2kuwLI"

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+correctJWT)
		req.Header.Set("X-Real-Ip", testIP)

		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(Equal("1"))
	})

	It("should allow to make normal response when X-Real-IP header is not provided, but RemoteAddr is", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImV4cCI6NDEwMjQ0NDgwMH0.4svEQGu2aNhKryzBiAOTTaygvomLzPrA8yEe_2kuwLI"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("X-Real-Ip", testIP)

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
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+unsupportedSigningMethodToken)
		req.Header.Set("X-Real-Ip", testIP)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Invalid authentication token."}}))
	})

	It("shoud return status unauthorized when token has wrong signature", func() {
		wrongSignatureToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImV4cCI6NDEwMjQ0NDgwMH0.NOFvlZPrNKfPTJ1yip6wzOMYuemmwo2U-g1m7KbSLc4"

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+wrongSignatureToken)
		req.Header.Set("X-Real-Ip", testIP)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Invalid authentication token."}}))
	})

	It("should return status unauthorized when token is expired", func() {
		expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImV4cCI6OTQ2Njg0ODAwfQ.c3NpXCy2wmB_-4oSys1FteHvkH-ikuZnSKy0aOfSji4"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		req.Header.Set("X-Real-Ip", testIP)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token has expired."}}))
	})

	// Test the case when token does not contain 'exp' field. This is incorrect and should be rejected.
	It("should return status unauthorized when token does not contain `exp` field", func() {
		expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSJ9.G_l6hjQr_CywF62Jnw2J9rb9CsNbEP0o5WG_kgSGH40"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		req.Header.Set("X-Real-Ip", testIP)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token has expired."}}))
	})

	// Test the case when token does not contain 'userID' field. This is incorrect and should be rejected.
	It("should return status unauthorized when token does not cotain `userID` field", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhbGxvd2VkSVAiOiIxOTIuMC4yLjEiLCJleHAiOjQxMDI0NDQ4MDB9.K7utB4yGFK11Bb78TMbk00cl9rzK1TjV2pdHsiUPW_I"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("X-Real-Ip", testIP)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token does not contain required data."}}))
	})

	It("should return status unauthorized when request is sent with X-Real-IP header that does not match `allowedIP`", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjEyNy4wLjAuMSIsImV4cCI6NDEwMjQ0NDgwMH0.Jm9nRlZgK1p-9bPaHBTD_E_31_ZmCdleL7uapCW5ZuI"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("X-Real-Ip", "1.2.3.4:1234")
		req.RemoteAddr = "" // just to make sure it checks X-Real-IP

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token is not allowed to be used from this IP."}}))
	})

	It("should return status unauthorized when request is sent without X-Real-IP header but with RemoteAddr that does not match `allowedIp`", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjEyNy4wLjAuMSIsImV4cCI6NDEwMjQ0NDgwMH0.Jm9nRlZgK1p-9bPaHBTD_E_31_ZmCdleL7uapCW5ZuI"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.RemoteAddr = "1.2.3.4:1234"

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token is not allowed to be used from this IP."}}))
	})

	// Not sure if this can ever happen
	It("should return status unauthorized when both X-Real-IP header and RemoteAddr are not provided", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImV4cCI6NDEwMjQ0NDgwMH0.4svEQGu2aNhKryzBiAOTTaygvomLzPrA8yEe_2kuwLI"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.RemoteAddr = ""

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Malformed request: no client IP was provided."}}))
	})

})
