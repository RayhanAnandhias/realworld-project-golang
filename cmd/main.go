package main

import (
	"log"
	"net/http"

	"github.com/RayhanAnandhias/realworld-project-golang/configs"
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/controllers"
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine

	TagController      controllers.TagController
	TagRouteController routes.TagRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController
)

func init() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	configs.ConnectDB(&config)

	TagController = controllers.NewTagController(configs.DB)
	TagRouteController = routes.NewTagRouteController(TagController)

	UserController = controllers.NewUserController(configs.DB)
	UserRouteController = routes.NewUserRouteController(UserController)

	server = gin.Default()
}

func main() {
	config, err := configs.LoadConfig("../")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Realworld Project"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	TagRouteController.TagRoute(router)
	UserRouteController.UserRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
