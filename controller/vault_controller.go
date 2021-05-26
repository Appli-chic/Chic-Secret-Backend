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

type VaultController struct {
	vaultService *service.VaultService
}

func NewVaultController() *VaultController {
	vaultController := new(VaultController)
	vaultController.vaultService = new(service.VaultService)
	return vaultController
}

func (v *VaultController) SaveVaults(c *gin.Context) {
	saveVaultsForm := validator2.SaveVaultsForm{}
	if err := c.ShouldBindJSON(&saveVaultsForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(saveVaultsForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the saveVaultsForm
	for _, vault := range saveVaultsForm.Vaults {
		currentCategory, err := v.vaultService.GetVault(vault.ID)

		if err == nil && currentCategory.UpdatedAt.Unix() < vault.UpdatedAt.Unix() {
			saveError := v.vaultService.Save(&vault)

			if saveError != nil {
				_ = fmt.Sprintf("VaultController->SaveVaults: %s", saveError.Error())
			}
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}

func (v *VaultController) GetVaults(c *gin.Context) {
	var vaults []model.Vault
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

		vaults, err = v.vaultService.GetVaultsToSynchronize(user.ID, lastSynchro)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Retrieve all the vaults
		vaults, err = v.vaultService.GetUserVaults(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"vaults": vaults,
	})
}
