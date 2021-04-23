package main

import (
	"applichic.com/chic_secret/controller"
	"applichic.com/chic_secret/util"
	"github.com/gin-gonic/gin"
)

// InitRouter Creates all the routes
func InitRouter() *gin.Engine {
	router := gin.Default()

	// Create controllers
	authController := controller.NewAuthController()

	api := router.Group("/api")
	{
		// Auth routes
		api.POST("/auth/ask_code", authController.AskCodeToLogin)
		api.POST("/auth/login", authController.Login)
		api.POST("/auth/refresh", authController.RefreshAccessToken)

		// Need to be logged in routes
		loggedInGroup := api.Group("/")
		loggedInGroup.Use(util.AuthenticationRequired())
		{
			loggedInGroup.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
		}
	}

	return router
}
