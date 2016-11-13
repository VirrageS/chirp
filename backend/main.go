package main

/*
	TODO:
	  - generate a good secret key
	    (useful: http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key,
	    	     https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/)
*/

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/api"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/middleware"
	"github.com/VirrageS/chirp/backend/service"
)

func init() {
	log.SetOutput(os.Stderr) // setup logrus logging library

}

func main() {
	router := createServer()
	router.Run(":8080")
}

// TODO: Move all setup to another package or file
func createServer() *gin.Engine {
	// setup Database
	DBConnection := database.NewDatabaseConnection()
	userDB := database.NewUserDB(DBConnection)

	service := service.NewService(userDB)
	api := api.NewAPI(service)

	return setupRouter(api)
}

func setupRouter(api *api.API) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(*setupCORS()))
	router.Use(middleware.ErrorHandler())

	contentTypeChecker := middleware.ContentTypeChecker()

	authorizedRoutes := router.Group("/", middleware.TokenAuthenticator)
	{
		tweets := authorizedRoutes.Group("tweets")
		tweets.GET("", api.GetTweets)
		tweets.POST("", contentTypeChecker, api.PostTweet)
		tweets.GET("/:id", api.GetTweet)
		tweets.DELETE("/:id", api.DeleteTweet)

		homeFeed := authorizedRoutes.Group("home_feed")
		homeFeed.GET("", api.HomeFeed)

		users := authorizedRoutes.Group("users")
		users.GET("", api.GetUsers)
		users.GET("/:id", api.GetUser)
	}

	accounts := router.Group("")
	{
		accounts.POST("/signup", contentTypeChecker, api.RegisterUser)
		accounts.POST("/login", contentTypeChecker, api.LoginUser)
	}

	return router
}

func setupCORS() *cors.Config {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	return &config
}
