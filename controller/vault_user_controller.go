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

type VaultUserController struct {
	vaultUserService *service.VaultUserService
}

func NewVaultUserController() *VaultUserController {
	vaultUserController := new(VaultUserController)
	vaultUserController.vaultUserService = new(service.VaultUserService)
	return vaultUserController
}

// SaveVaultUser Save the vault users to synchronize in the database
func (vaultUserController *VaultUserController) SaveVaultUser(c *gin.Context) {
	vaultUserForm := validator2.SaveVaultUserForm{}
	if err := c.ShouldBindJSON(&vaultUserForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(vaultUserForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the vault users
	for _, vaultUser := range vaultUserForm.VaultUsers {
		currentVaultUser, err := vaultUserController.vaultUserService.GetVaultUser(vaultUser.VaultID, vaultUser.UserID)

		if err == nil && currentVaultUser.UpdatedAt.Unix() < vaultUser.UpdatedAt.Unix() {
			saveError := vaultUserController.vaultUserService.Save(&vaultUser)

			if saveError != nil {
				_ = fmt.Sprintf("VaultUserController->SaveVaultUser: %s", saveError.Error())
			}
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}

// GetVaultUsers Retrieve the vault users to synchronize with the user's device
func (vaultUserController *VaultUserController) GetVaultUsers(c *gin.Context) {
	var vaultUsers []model.VaultUser
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
		// Retrieve the vault users that changed
		lastSynchro, err := time.Parse(layout, lastSynchroString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		vaultUsers, err = vaultUserController.vaultUserService.GetVaultUsersToSynchronize(user.ID, lastSynchro)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Retrieve all the vault users
		vaultUsers, err = vaultUserController.vaultUserService.GetVaultUsers(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"vault_users": vaultUsers,
	})
}
