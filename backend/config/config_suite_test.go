package config

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// NOTE: DO NOT CHANGE SPACES TO TABS HERE! - it will break loading
var content = []byte(`
server_defaults: &server_defaults
  secret_key: "just a random secret string"
  auth_token_validity_period: 15m
  refresh_token_validity_period: 24h
  random_password_length: 128

database_defaults: &database_defaults
  username: "postgres"
  password: "postgres"
  host: "localhost"
  port: "5432"

redis_defaults: &redis_defaults
  password: "pass"
  host: "localhost"
  port: "6379"
  db: 0
  expiration_time: 1m

authorization_google_defaults: &authorization_google_defaults
  client_id: "248788072320-cgap8rml3940qugk1u1i1c77onukabnn.apps.googleusercontent.com"
  client_secret: "V4lFJ6QWLP117cNN9Y2O3wGj"
  callback_uri: "http://localhost:3000/login/google/callback"
  auth_url: "https://accounts.google.com/o/oauth2/auth"
  token_url: "https://accounts.google.com/o/oauth2/token"

elasticsearch_defaults: &elasticsearch_defaults
  username: "elastic"
  password: "changeme"
  host: "localhost"
  port: "9200"

defaults: &defaults
  <<: *server_defaults
  database:
    <<: *database_defaults
  redis:
    <<: *redis_defaults
  authorization_google:
    <<: *authorization_google_defaults
  elasticsearch:
    <<: *elasticsearch_defaults

development:
  <<: *defaults

production:
  <<: *defaults
  database:
    <<: *database_defaults
    host: "database"
  redis:
    <<: *redis_defaults
    host: "cache"
  authorization_google:
    <<: *authorization_google_defaults
    callback_uri: "http://frontend.show/login/google/callback"
  elasticsearch:
    <<: *elasticsearch_defaults
    host: "elasticsearch"

test:
  <<: *defaults
  secret_key: "secret"
  random_password_length: 32
  database:
    <<: *database_defaults
    port: "5433"
  redis:
    <<: *redis_defaults
    port: "6380"
    expiration_time: 50ms
`)

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config")
}
