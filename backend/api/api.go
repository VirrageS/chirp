package api

import (
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/service"

	"golang.org/x/oauth2"
)

// Struct that implements APIProvider
type API struct {
	// logger?
	service      service.ServiceProvider
	googleOAuth2 oauth2.Config
}

// Constructs an API object that uses given ServiceProvider.
func NewAPI(
	service service.ServiceProvider,
	authorizationGoogleConfig config.AuthorizationGoogleConfigurationProvider,
) APIProvider {
	googleOAuth2 := oauth2.Config{
		ClientID:     authorizationGoogleConfig.GetClientID(),
		ClientSecret: authorizationGoogleConfig.GetClientSecret(),
		RedirectURL:  authorizationGoogleConfig.GetCallbackURI(),
		Scopes:       []string{"email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizationGoogleConfig.GetAuthURL(),
			TokenURL: authorizationGoogleConfig.GetTokenURL(),
		},
	}

	return &API{
		service:      service,
		googleOAuth2: googleOAuth2,
	}
}
