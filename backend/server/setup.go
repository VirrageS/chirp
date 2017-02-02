package server

import (
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/middleware"
	"github.com/VirrageS/chirp/backend/password"
	"github.com/VirrageS/chirp/backend/service"
	"github.com/VirrageS/chirp/backend/storage"
	"github.com/VirrageS/chirp/backend/token"
)

// New creates a new server.
func New() *gin.Engine {
	conf := config.New()
	if conf == nil {
		panic("Failed to get config.")
	}

	storage := storage.New(conf.Postgres, conf.Redis, conf.Elasticsearch)
	passwordManager := password.NewBcryptManager(conf.Password)
	services := service.New(storage, passwordManager)

	tokenManager := token.NewManager(conf.Token)
	apis := api.New(services, tokenManager, conf.AuthorizationGoogle)

	return setupRouter(apis, tokenManager)
}

func setupRouter(api api.APIProvider, tokenManager token.Manager) *gin.Engine {
	corsHandler := newCorsHandler()
	contentTypeChecker := middleware.ContentTypeChecker()
	authenticator := middleware.TokenAuthenticator(tokenManager)
	errorHandler := middleware.ErrorHandler()

	router := gin.Default()
	router.Use(corsHandler)
	router.Use(errorHandler)

	authorizedRoutes := router.Group("/", authenticator)
	{
		tweets := authorizedRoutes.Group("tweets")
		tweets.POST("", contentTypeChecker, api.PostTweet)
		tweets.GET("/:id", api.GetTweet)
		tweets.DELETE("/:id", api.DeleteTweet)
		tweets.POST("/:id/like", api.LikeTweet)
		tweets.POST("/:id/unlike", api.UnlikeTweet)

		feed := authorizedRoutes.Group("feed")
		feed.GET("", api.Feed)

		users := authorizedRoutes.Group("users")
		users.GET("/:id", api.GetUser)
		users.POST(":id/follow", api.FollowUser)
		users.POST(":id/unfollow", api.UnfollowUser)
		users.GET(":id/followers", api.UserFollowers)
		users.GET(":id/followees", api.UserFollowees)
		users.GET(":id/tweets", api.UserTweets)

		search := authorizedRoutes.Group("search")
		search.GET("", api.Search)
	}

	auth := router.Group("")
	{
		auth.POST("/signup", contentTypeChecker, api.RegisterUser)
		auth.POST("/login", contentTypeChecker, api.LoginUser)
		auth.POST("/token", contentTypeChecker, api.RefreshAuthToken)
		auth.GET("/authorize/google", api.GetGoogleAuthorizationURL)
		auth.POST("/login/google", api.CreateOrLoginUserWithGoogle)
	}

	return router
}

func newCorsHandler() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	return cors.New(config)
}
