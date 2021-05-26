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

type CustomFieldController struct {
	customFieldService *service.CustomFieldService
}

func NewCustomFieldController() *CustomFieldController {
	customFieldController := new(CustomFieldController)
	customFieldController.customFieldService = new(service.CustomFieldService)
	return customFieldController
}

func (customFieldController *CustomFieldController) SaveCustomFields(c *gin.Context) {
	customFieldForm := validator2.SaveCustomFieldForm{}
	if err := c.ShouldBindJSON(&customFieldForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(customFieldForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the customFieldForm
	for _, customField := range customFieldForm.CustomFields {
		currentCustomField, err := customFieldController.customFieldService.GetCustomField(customField.ID)

		if err == nil && currentCustomField.UpdatedAt.Unix() < customField.UpdatedAt.Unix() {
			saveError := customFieldController.customFieldService.Save(&customField)

			if saveError != nil {
				_ = fmt.Sprintf("CustomFieldController->SaveCustomFields: %s", saveError.Error())
			}
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}

func (customFieldController *CustomFieldController) GetCustomFields(c *gin.Context) {
	var customFields []model.CustomField
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
		// Retrieve the categories that changed
		lastSynchro, err := time.Parse(layout, lastSynchroString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		customFields, err = customFieldController.customFieldService.GetCustomFieldsToSynchronize(user.ID, lastSynchro)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Retrieve all the custom fields
		customFields, err = customFieldController.customFieldService.GetCustomFieldsFromVault(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"custom_fields": customFields,
	})
}
