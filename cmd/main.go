package main

import (
	"log"
	"net/http"

	"github.com/RayhanAnandhias/realworld-project-golang/configs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
)

func init() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	configs.ConnectDB(&config)

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

	log.Fatal(server.Run(":" + config.ServerPort))
}
