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

type CategoryController struct {
	categoryService *service.CategoryService
}

func NewCategoryController() *CategoryController {
	categoryController := new(CategoryController)
	categoryController.categoryService = new(service.CategoryService)
	return categoryController
}

// SaveCategories Save the categories to synchronize in the database
func (categoryController *CategoryController) SaveCategories(c *gin.Context) {
	categoriesForm := validator2.SaveCategoriesForm{}
	if err := c.ShouldBindJSON(&categoriesForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(categoriesForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the categoriesForm
	for _, category := range categoriesForm.Categories {
		currentCategory, err := categoryController.categoryService.GetCategory(category.ID)

		if err == nil && currentCategory.UpdatedAt.Unix() < category.UpdatedAt.Unix() {
			saveError := categoryController.categoryService.Save(&category)

			if saveError != nil {
				_ = fmt.Sprintf("CategoryController->SaveCategories: %s", saveError.Error())
			}
		}
	}

	c.JSONP(http.StatusOK, gin.H{})
}

// GetCategories Retrieve the categories to synchronize with the user's device
func (categoryController *CategoryController) GetCategories(c *gin.Context) {
	var categories []model.Category
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

		categories, err = categoryController.categoryService.GetCategoriesToSynchronize(user.ID, lastSynchro)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Retrieve all the categories
		categories, err = categoryController.categoryService.GetCategoriesFromVault(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"categories": categories,
	})
}
