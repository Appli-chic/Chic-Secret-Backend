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
	entryController := controller.NewEntryController()
	customFieldController := controller.NewCustomFieldController()
	tagController := controller.NewTagController()
	entryTagController := controller.NewEntryTagController()
	vaultUserController := controller.NewVaultUserController()

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
			loggedInGroup.POST("/user", userController.SaveUser)
			loggedInGroup.GET("/user", userController.FetchUser)
			loggedInGroup.GET("/user/email", userController.FetchUserByEmail)
			loggedInGroup.GET("/users", userController.GetUsers)

			// Vaults
			loggedInGroup.POST("/vaults", vaultController.SaveVaults)
			loggedInGroup.GET("/vaults", vaultController.GetVaults)

			// Categories
			loggedInGroup.POST("/categories", categoryController.SaveCategories)
			loggedInGroup.GET("/categories", categoryController.GetCategories)

			// Entries
			loggedInGroup.POST("/entries", entryController.SaveEntries)
			loggedInGroup.GET("/entries", entryController.GetEntries)

			// Custom Fields
			loggedInGroup.POST("/custom-fields", customFieldController.SaveCustomFields)
			loggedInGroup.GET("/custom-fields", customFieldController.GetCustomFields)

			// Tags
			loggedInGroup.POST("/tags", tagController.SaveTags)
			loggedInGroup.GET("/tags", tagController.GetTags)

			// Entry Tags
			loggedInGroup.POST("/entry-tags", entryTagController.SaveEntryTags)
			loggedInGroup.GET("/entry-tags", entryTagController.GetEntryTags)

			// Vault Users
			loggedInGroup.POST("/vault-users", vaultUserController.SaveVaultUser)
			loggedInGroup.GET("/vault-users", vaultUserController.GetVaultUsers)
		}
	}

	return router
}
