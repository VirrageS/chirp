package server

import (
	"github.com/gin-gonic/gin"

	"github.com/VirrageS/chirp/backend/api"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/password"
	"github.com/VirrageS/chirp/backend/service"
	"github.com/VirrageS/chirp/backend/storage"
	"github.com/VirrageS/chirp/backend/token"
)

// FakeServer exports all fields which can be necessary for tests.
type FakeServer struct {
	Server       *gin.Engine
	TokenManager token.Manager
	Storage      *storage.FakeStorage
}

// NewFakeServer creates a new fake server.
func NewFakeServer() *FakeServer {
	conf := config.New()
	if conf == nil {
		panic("Failed to get config.")
	}

	fakeStorage := storage.NewFakeStorage(conf.Postgres)
	passwordManager := password.NewBcryptManager(conf.Password)
	services := service.New(fakeStorage.Storage, passwordManager)

	tokenManager := token.NewManager(conf.Token)
	apis := api.New(services, tokenManager, conf.AuthorizationGoogle)

	return &FakeServer{
		Server:       setupRouter(apis, tokenManager),
		TokenManager: tokenManager,
		Storage:      fakeStorage,
	}
}
