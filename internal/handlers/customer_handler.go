package handlers

import (
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"

	"github.com/gin-gonic/gin"
)

func ListCustomers(c *gin.Context) {
	var customers []models.Customer
	
	// Basic search
	search := c.Query("q")
	query := configs.DB.Model(&models.Customer{})
	if search != "" {
		query = query.Where("name LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	query.Order("created_at desc").Find(&customers)

	c.HTML(http.StatusOK, "customers.html", gin.H{
		"Customers": customers,
		"Search":    search,
	})
}

func CreateCustomer(c *gin.Context) {
	customer := models.Customer{
		Name:        c.PostForm("name"),
		Phone:       c.PostForm("phone"),
		Address:     c.PostForm("address"),
		IsActive:    c.PostForm("is_active") == "on",
		Notes:       c.PostForm("notes"),
	}

	if err := configs.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/customers")
}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer
	if err := configs.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	customer.Name = c.PostForm("name")
	customer.Phone = c.PostForm("phone")
	customer.Address = c.PostForm("address")
	customer.IsActive = c.PostForm("is_active") == "on"
	customer.Notes = c.PostForm("notes")

	if err := configs.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/customers")
}

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := configs.DB.Delete(&models.Customer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/customers")
}
