package controller

import (
	"applichic.com/chic_secret/model"
	"applichic.com/chic_secret/service"
	"applichic.com/chic_secret/util"
	validator2 "applichic.com/chic_secret/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type UserController struct {
	userService        *service.UserService
	tagService         *service.TagService
	entryTagService    *service.EntryTagService
	customFieldService *service.CustomFieldService
	entryService       *service.EntryService
	categoryService    *service.CategoryService
	vaultUserService   *service.VaultUserService
	vaultService       *service.VaultService
	tokenService       *service.TokenService
	loginTokenService  *service.LoginTokenService
}

func NewUserController() *UserController {
	userController := new(UserController)
	userController.userService = new(service.UserService)
	userController.tagService = new(service.TagService)
	userController.customFieldService = new(service.CustomFieldService)
	userController.entryService = new(service.EntryService)
	userController.categoryService = new(service.CategoryService)
	userController.vaultUserService = new(service.VaultUserService)
	userController.vaultService = new(service.VaultService)
	userController.tokenService = new(service.TokenService)
	userController.loginTokenService = new(service.LoginTokenService)
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

// FetchUserByEmail Fetch user by email
func (u *UserController) FetchUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := u.userService.FetchUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	currentUser, err := u.userService.FetchUserById(saveUserForm.User.ID)

	if err == nil && currentUser.UpdatedAt.Unix() < saveUserForm.User.UpdatedAt.Unix() {
		saveError := u.userService.Save(&saveUserForm.User)

		if saveError != nil {
			_ = fmt.Sprintf("UserController->Save: %s", saveError.Error())
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}

func (u *UserController) GetUsers(c *gin.Context) {
	var users []model.User
	layout := "2006-01-02T15:04:05Z"
	lastSynchroString := c.Query("LastSynchro")

	// Check if the user exists
	user, err := util.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	if lastSynchroString != "" && lastSynchroString != "null" {
		// Retrieve the vaults that changed
		lastSynchro, err := time.Parse(layout, lastSynchroString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		users, err = u.userService.GetUsersToSynchronize(user.ID, lastSynchro)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Retrieve all the users
		users, err = u.userService.GetUsersLinkedToUser(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"users": users,
	})
}

func (u *UserController) Delete(c *gin.Context) {
	// Check if the user exists
	user, err := util.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	u.entryTagService.DeleteFromUser(user.ID)
	u.tagService.DeleteFromUser(user.ID)
	u.customFieldService.DeleteFromUser(user.ID)
	u.entryService.DeleteFromUser(user.ID)
	u.categoryService.DeleteFromUser(user.ID)
	u.vaultUserService.DeleteFromUser(user.ID)
	u.vaultService.DeleteFromUser(user.ID)
	u.tokenService.DeleteFromUser(user.ID)
	u.loginTokenService.DeleteFromUser(user.ID)
	u.userService.DeleteFromUser(user.ID)

	c.JSONP(http.StatusOK, gin.H{})
}
