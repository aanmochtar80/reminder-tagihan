package handlers

import (
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowDashboard(c *gin.Context) {
	session := sessions.Default(c)
	flashMsg := session.Get("flash")
	if flashMsg != nil {
		session.Delete("flash")
		session.Save()
	}

	var totalCustomers int64
	configs.DB.Model(&models.Customer{}).Count(&totalCustomers)

	// In a real app, calculate these based on the current month/day
	stats := map[string]interface{}{
		"TotalCustomers":     totalCustomers,
		"TotalInvoicesMonth": "0", // Calculate later
		"DueToday":           0,   // Calculate later
		"Overdue":            0,   // Calculate later
	}

	var recentInvoices []models.Invoice
	// Eager load customer for recent invoices
	configs.DB.Preload("Customer").Order("created_at desc").Limit(5).Find(&recentInvoices)

	// Format data for template
	type InvoiceDisplay struct {
		CustomerName string
		Status       string
		Amount       float64
	}
	var displays []InvoiceDisplay
	for _, inv := range recentInvoices {
		displays = append(displays, InvoiceDisplay{
			CustomerName: inv.Customer.Name,
			Status:       inv.Status,
			Amount:       inv.Amount,
		})
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"FlashMessage":   flashMsg,
		"Stats":          stats,
		"RecentInvoices": displays,
	})
}
