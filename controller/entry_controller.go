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

type EntryController struct {
	entryService *service.EntryService
}

func NewEntryController() *EntryController {
	entryController := new(EntryController)
	entryController.entryService = new(service.EntryService)
	return entryController
}

// SaveEntries Save the entries to synchronize in the database
func (e *EntryController) SaveEntries(c *gin.Context) {
	entriesForm := validator2.SaveEntriesForm{}
	if err := c.ShouldBindJSON(&entriesForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(entriesForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the entriesForm
	for _, entry := range entriesForm.Entries {
		currentEntry, err := e.entryService.GetEntry(entry.ID)

		if err == nil && currentEntry.UpdatedAt.Unix() < entry.UpdatedAt.Unix() {
			saveError := e.entryService.Save(&entry)

			if saveError != nil {
				_ = fmt.Sprintf("EntryController->SaveEntries: %s", saveError.Error())
			}
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}

// GetEntries Retrieve the entries to synchronize with the user's device
func (e *EntryController) GetEntries(c *gin.Context) {
	var entries []model.Entry
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
		// Retrieve the entries that changed
		lastSynchro, err := time.Parse(layout, lastSynchroString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		entries, err = e.entryService.GetEntriesToSynchronize(user.ID, lastSynchro)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Retrieve all the entries
		entries, err = e.entryService.GetEntriesFromVault(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"entries": entries,
	})
}
