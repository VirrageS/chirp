package api

import (
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/service"

	"github.com/VirrageS/chirp/backend/token"
	"golang.org/x/oauth2"
)

// Struct that implements APIProvider
type API struct {
	service      service.ServiceProvider
	tokenManager token.Manager
	googleOAuth2 oauth2.Config
}

// Constructs an API object that uses given ServiceProvider.
func New(
	service service.ServiceProvider,
	tokenManager token.Manager,
	authorizationGoogleConfig config.AuthorizationGoogleConfigProvider,
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
		tokenManager: tokenManager,
		googleOAuth2: googleOAuth2,
	}
}
