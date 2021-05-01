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

type EntryTagController struct {
	entryTagService *service.EntryTagService
}

func NewEntryTagController() *EntryTagController {
	entryTagFieldController := new(EntryTagController)
	entryTagFieldController.entryTagService = new(service.EntryTagService)
	return entryTagFieldController
}

func (entryTagController *EntryTagController) SaveEntryTags(c *gin.Context) {
	entryTagForm := validator2.SaveEntryTagForm{}
	if err := c.ShouldBindJSON(&entryTagForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(entryTagForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the entry tags
	for _, entryTag := range entryTagForm.EntryTags {
		saveError := entryTagController.entryTagService.Save(&entryTag)

		if saveError != nil {
			_ = fmt.Sprintf("EntryTagController->SaveEntryTags: %s", saveError.Error())
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}

func (entryTagController *EntryTagController) GetEntryTags(c *gin.Context) {
	var entryTags []model.EntryTag
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
		// Retrieve the tags that changed
		lastSynchro, err := time.Parse(layout, lastSynchroString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		entryTags, err = entryTagController.entryTagService.GetEntryTagsToSynchronize(user.ID, lastSynchro)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Retrieve all the custom fields
		entryTags, err = entryTagController.entryTagService.GetEntryTagsFromVault(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"entry_tags": entryTags,
	})
}
