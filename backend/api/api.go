package api

import "github.com/VirrageS/chirp/backend/service"

// Struct that implements APIProvider
type API struct {
	// logger?
	Service service.ServiceProvider
}

// Constructs an API object that uses given ServiceProvider.
func NewAPI(service service.ServiceProvider) APIProvider {
	return &API{service}
}
