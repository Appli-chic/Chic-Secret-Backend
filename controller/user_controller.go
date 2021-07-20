package controller

import (
	"applichic.com/chic_secret/service"
	"applichic.com/chic_secret/util"
	validator2 "applichic.com/chic_secret/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	userController := new(UserController)
	userController.userService = new(service.UserService)
	return userController
}

// FetchUser Fetch user's data
func (u *UserController) FetchUser(c *gin.Context) {
	user, err := util.GetUserFromToken(c)

	// Check if the user exists
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Send the user information
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// SaveUser Save and Update the user
func (u *UserController) SaveUser(c *gin.Context) {
	saveUserForm := validator2.SaveUserForm{}
	if err := c.ShouldBindJSON(&saveUserForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(saveUserForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the saveUserForm
	currentCategory, err := u.userService.FetchUserById(saveUserForm.User.ID)

	if err == nil && currentCategory.UpdatedAt.Unix() < saveUserForm.User.UpdatedAt.Unix() {
		saveError := u.userService.Save(&saveUserForm.User)

		if saveError != nil {
			_ = fmt.Sprintf("UserController->Save: %s", saveError.Error())
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}
