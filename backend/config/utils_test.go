package config

import (
	"bytes"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("Utils", func() {
	var (
		config *generalConfig
	)

	BeforeEach(func() {
		v := viper.New()
		v.SetConfigType("yaml")
		v.ReadConfig(bytes.NewBuffer(content))
		config = &generalConfig{v.Sub("development")}
	})

	It("should not return nil when reading content", func() {
		Expect(config).NotTo(BeNil())
	})

	It("should return proper values for Token config provider", func() {
		var token TokenConfigProvider = config.getServerConfig()
		Expect(token.GetSecretKey()).To(Equal([]byte("just a random secret string")))
		Expect(token.GetAuthTokenValidityPeriod()).To(Equal(time.Duration(15) * time.Minute))
		Expect(token.GetRefreshTokenValidityPeriod()).To(Equal(time.Duration(24) * time.Hour))
	})

	It("should return proper values for Password config provider", func() {
		var password PasswordConfigProvider = config.getServerConfig()
		Expect(password.GetRandomPasswordLength()).To(Equal(128))
	})

	It("should return proper values for Database config provider", func() {
		var database DatabaseConfigProvider = config.getDatabaseConfig()
		Expect(database.GetUsername()).To(Equal("postgres"))
		Expect(database.GetPassword()).To(Equal("postgres"))
		Expect(database.GetHost()).To(Equal("localhost"))
		Expect(database.GetPort()).To(Equal("5432"))
	})

	It("should return proper values for Redis config provider", func() {
		var redis RedisConfigProvider = config.getRedisCacheConfig()
		Expect(redis.GetPassword()).To(Equal("pass"))
		Expect(redis.GetHost()).To(Equal("localhost"))
		Expect(redis.GetPort()).To(Equal("6379"))
		Expect(redis.GetExpirationTime()).To(Equal(time.Duration(1) * time.Minute))
		Expect(redis.GetDB()).To(Equal(0))
	})

	It("should return proper values for Google Authorization config provider", func() {
		var google AuthorizationGoogleConfigProvider = config.getAuthorizationGoogleConfig()
		Expect(google.GetClientID()).To(Equal("248788072320-cgap8rml3940qugk1u1i1c77onukabnn.apps.googleusercontent.com"))
		Expect(google.GetClientSecret()).To(Equal("V4lFJ6QWLP117cNN9Y2O3wGj"))
		Expect(google.GetCallbackURI()).To(Equal("http://localhost:3000/login/google/callback"))
		Expect(google.GetAuthURL()).To(Equal("https://accounts.google.com/o/oauth2/auth"))
		Expect(google.GetTokenURL()).To(Equal("https://accounts.google.com/o/oauth2/token"))
	})

	It("should return proper values for Elasticsearch config provider", func() {
		var elastic ElasticsearchConfigProvider = config.getElasticsearchConfig()
		Expect(elastic.GetUsername()).To(Equal("elastic"))
		Expect(elastic.GetPassword()).To(Equal("changeme"))
		Expect(elastic.GetHost()).To(Equal("localhost"))
		Expect(elastic.GetPort()).To(Equal("9200"))
	})
})
