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
	userController := controller.NewUserController()
	vaultController := controller.NewVaultController()
	categoryController := controller.NewCategoryController()

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
			// User
			loggedInGroup.GET("/user", userController.FetchUser)

			// Vaults
			loggedInGroup.POST("/vaults", vaultController.SaveVaults)
			loggedInGroup.GET("/vaults", vaultController.GetVaults)

			// Categories
			loggedInGroup.POST("/categories", categoryController.SaveCategories)
			loggedInGroup.GET("/categories", categoryController.GetCategories)
		}
	}

	return router
}
