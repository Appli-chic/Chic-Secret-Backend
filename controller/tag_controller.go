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

type TagController struct {
	tagService *service.TagService
}

func NewTagController() *TagController {
	customFieldController := new(TagController)
	customFieldController.tagService = new(service.TagService)
	return customFieldController
}

func (tagController *TagController) SaveTags(c *gin.Context) {
	tagForm := validator2.SaveTagForm{}
	if err := c.ShouldBindJSON(&tagForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(tagForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the tags
	for _, tag := range tagForm.Tags {
		currentCategory, err := tagController.tagService.GetTag(tag.ID)

		if err == nil && currentCategory.UpdatedAt.Unix() < tag.UpdatedAt.Unix() {
			saveError := tagController.tagService.Save(&tag)

			if saveError != nil {
				_ = fmt.Sprintf("TagController->SaveTags: %s", saveError.Error())
			}
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}

func (tagController *TagController) GetTags(c *gin.Context) {
	var tags []model.Tag
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

		tags, err = tagController.tagService.GetTagsToSynchronize(user.ID, lastSynchro)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Retrieve all the custom fields
		tags, err = tagController.tagService.GetTagsFromVault(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"tags": tags,
	})
}
