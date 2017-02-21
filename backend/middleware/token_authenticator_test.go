package middleware

import (
	"net/http"
	"net/http/httptest"

	"encoding/json"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/token"
)

// From Go documentation of httptest.NewRequest():
// 192.0.2.0/24 is "TEST-NET" in RFC 5737 for use solely in
// documentation and example source code and should not be
// used publicly.
const (
	testIP    = "192.0.2.1"
	testAgent = "test/1.0"
)

var _ = Describe("TokenAuthenticator", func() {
	var (
		router *gin.Engine
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)

		conf := config.New()
		tokenManager := token.NewManager(conf.Token)
		router = gin.New()
		router.Use(ErrorHandler())
		router.Use(TokenAuthenticator(tokenManager))
	})

	It("should allow to make normal response when jwt token is okay", func() {
		correctJWT := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6NDEwMjQ0NDgwMH0.OO-BsNijn7kIG1G2QUhNPibdTZF2Q9ezpHSQ6JJ_AwE"

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+correctJWT)
		req.Header.Set("X-Real-Ip", testIP)
		req.Header.Set("User-Agent", testAgent)

		router.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(Equal("1"))
	})

	It("should allow to make normal response when X-Real-IP header is not provided, but RemoteAddr is", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6NDEwMjQ0NDgwMH0.OO-BsNijn7kIG1G2QUhNPibdTZF2Q9ezpHSQ6JJ_AwE"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("X-Real-Ip", testIP)
		req.Header.Set("User-Agent", testAgent)

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
		req.Header.Set("User-Agent", testAgent)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Invalid authentication token."}}))
	})

	It("shoud return status unauthorized when token has wrong signature", func() {
		wrongSignatureToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6NDEwMjQ0NDgwMH0.uFkBc1pQDK_YL8y1tbllqkrmnjk3gG_h8Bi_B45LFdY"

		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+wrongSignatureToken)
		req.Header.Set("X-Real-Ip", testIP)
		req.Header.Set("User-Agent", testAgent)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Invalid authentication token."}}))
	})

	It("should return status unauthorized when token is expired", func() {
		expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6OTQ2Njg0ODAwfQ.CGy0aYslZYNIz_CI5n9XjAsxBkulotkaEfbsR11WM5s"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		req.Header.Set("X-Real-Ip", testIP)
		req.Header.Set("User-Agent", testAgent)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token has expired."}}))
	})

	// Test the case when token does not contain 'exp' field. This is incorrect and should be rejected.
	It("should return status unauthorized when token does not contain `exp` field", func() {
		expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSJ9.vxk2n-A27G72Miu9geKKVjxKF5vE4Gfc88jzyLDiEW0"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		req.Header.Set("X-Real-Ip", testIP)
		req.Header.Set("User-Agent", testAgent)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token has expired."}}))
	})

	// Test the case when token does not contain 'userID' field. This is incorrect and should be rejected.
	It("should return status unauthorized when token does not cotain `userID` field", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhbGxvd2VkSVAiOiIxOTIuMC4yLjEiLCJhbGxvd2VkVXNlckFnZW50IjoidGVzdC8xLjAiLCJleHAiOjQxMDI0NDQ4MDB9.uxEGCy4krGYUZ00TLojsIVYUUY5mpzJ9VDh9DhwMhpE"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("X-Real-Ip", testIP)
		req.Header.Set("User-Agent", testAgent)

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token does not contain required data."}}))
	})

	It("should return status unauthorized when request is sent with X-Real-IP header that does not match `allowedIP`", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjEyNy4wLjAuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6NDEwMjQ0NDgwMH0.AC_FQLzE660E7RWrq6BiRUloINb6Iaydqi6i3oK3E8M"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("X-Real-Ip", "1.2.3.4:1234")
		req.Header.Set("User-Agent", testAgent)
		req.RemoteAddr = "" // just to make sure it checks X-Real-IP

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token is not allowed to be used from this IP."}}))
	})

	It("should return status unauthorized when request is sent without X-Real-IP header but with RemoteAddr that does not match `allowedIp`", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjEyNy4wLjAuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6NDEwMjQ0NDgwMH0.AC_FQLzE660E7RWrq6BiRUloINb6Iaydqi6i3oK3E8M"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("User-Agent", testAgent)
		req.RemoteAddr = "1.2.3.4:1234"

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token is not allowed to be used from this IP."}}))
	})

	// Not sure if this can ever happen
	It("should return status unauthorized when both X-Real-IP header and RemoteAddr are not provided", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6NDEwMjQ0NDgwMH0.OO-BsNijn7kIG1G2QUhNPibdTZF2Q9ezpHSQ6JJ_AwE"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("User-Agent", testAgent)
		req.RemoteAddr = ""

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Malformed request: no client IP was provided."}}))
	})

	It("should return status unauthorized when User-Agent header is not provided.", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6NDEwMjQ0NDgwMH0.OO-BsNijn7kIG1G2QUhNPibdTZF2Q9ezpHSQ6JJ_AwE"
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
		Expect(response).To(Equal(errorResponse{[]string{"Malformed request: no User-Agent header."}}))
	})

	It("should return status unauthorized when User-Agent header is not provided.", func() {
		noUserIDToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VySUQiOjEsImFsbG93ZWRJUCI6IjE5Mi4wLjIuMSIsImFsbG93ZWRVc2VyQWdlbnQiOiJ0ZXN0LzEuMCIsImV4cCI6NDEwMjQ0NDgwMH0.OO-BsNijn7kIG1G2QUhNPibdTZF2Q9ezpHSQ6JJ_AwE"
		router.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "%d", c.MustGet("userID").(int64))
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+noUserIDToken)
		req.Header.Set("X-Real-Ip", testIP)
		req.Header.Set("User-Agent", "badAgent/1.0")

		router.ServeHTTP(w, req)

		var response errorResponse
		json.NewDecoder(w.Body).Decode(&response)

		Expect(w.Code).To(Equal(http.StatusUnauthorized))
		Expect(response).To(Equal(errorResponse{[]string{"Token is not allowed to be used from this User-Agent."}}))
	})
})
