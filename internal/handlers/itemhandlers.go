// itemhandlers.go

package handlers

import (
	"fmt"
	"go_project/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (repo *ItemRepository) GetAllItems(c *gin.Context) {
	var items []models.Item
	result := repo.db.Find(&items)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (repo *ItemRepository) GetItemByID(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	result := repo.db.First(&item, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (repo *ItemRepository) CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	result := repo.db.Where("name = ?", item.Category.Name).First(&category)

	if category.ID == 0 {
		repo.db.Create(category)
		fmt.Println("my category =", category)
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	item.Category = &category

	result = repo.db.Create(&item)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (repo *ItemRepository) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := repo.db.Model(&models.Item{}).Where("id = ?", id).Updates(&item)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	repo.db.Where("id=?", id).First(&item)
	c.JSON(http.StatusOK, item)
}

func (repo *ItemRepository) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	result := repo.db.Delete(&models.Item{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
