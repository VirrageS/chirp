package server

import (
	"database/sql"
	"os"

	"github.com/Sirupsen/logrus"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/middleware"
	"github.com/VirrageS/chirp/backend/service"
	"github.com/VirrageS/chirp/backend/token"
)

func init() {
	logrus.SetOutput(os.Stderr)
}

// Handles all dependencies and creates a new server.
// Takes a DB connection parameter in order to support test database.
func New(dbConnection *sql.DB) *gin.Engine {
	// service dependencies
	serverConfig := config.GetConfig()
	tokenManager := token.NewTokenManager(serverConfig)

	// api dependencies
	CORSConfig := setupCORS()

	db := database.NewDatabase(dbConnection)
	services := service.NewService(serverConfig, db, tokenManager)
	APIs := api.NewAPI(services)

	return setupRouter(APIs, tokenManager, CORSConfig)
}

// TODO: Maybe middlewares should also be dependencies
func setupRouter(api api.APIProvider, tokenManager token.TokenManagerProvider, corsConfig *cors.Config) *gin.Engine {
	CORSHandler := cors.New(*corsConfig)
	contentTypeChecker := middleware.ContentTypeChecker()
	authenticator := middleware.TokenAuthenticator(tokenManager)
	errorHandler := middleware.ErrorHandler()

	router := gin.Default()
	router.Use(CORSHandler)
	router.Use(errorHandler)

	authorizedRoutes := router.Group("/", authenticator)
	{
		tweets := authorizedRoutes.Group("tweets")
		tweets.GET("", api.GetTweets)
		tweets.POST("", contentTypeChecker, api.PostTweet)
		tweets.GET("/:id", api.GetTweet)
		tweets.DELETE("/:id", api.DeleteTweet)
		tweets.POST("/:id/like", api.LikeTweet)
		tweets.POST("/:id/unlike", api.UnlikeTweet)

		homeFeed := authorizedRoutes.Group("home_feed")
		homeFeed.GET("", api.HomeFeed)

		users := authorizedRoutes.Group("users")
		users.GET("", api.GetUsers)
		users.GET("/:id", api.GetUser)
		users.POST(":id/follow", api.FollowUser)
		users.POST(":id/unfollow", api.UnfollowUser)
	}

	auth := router.Group("")
	{
		auth.POST("/signup", contentTypeChecker, api.RegisterUser)
		auth.POST("/login", contentTypeChecker, api.LoginUser)
		auth.POST("/token", contentTypeChecker, api.RefreshAuthToken)
	}

	return router
}

func setupCORS() *cors.Config {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	return &config
}
