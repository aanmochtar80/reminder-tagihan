package handlers

import (
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"
	"time"

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

	now := time.Now()
	// Use UTC to match the timezone used by time.Parse("2006-01-02") so SQLite string comparisons work correctly
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := todayStart.Add(24 * time.Hour)
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var totalInvoicesMonth float64
	configs.DB.Model(&models.Invoice{}).
		Where("due_date >= ? AND due_date <= ?", firstOfMonth, lastOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalInvoicesMonth)

	var dueToday int64
	configs.DB.Model(&models.Invoice{}).
		Where("due_date >= ? AND due_date < ? AND status = ?", todayStart, todayEnd, "pending").
		Count(&dueToday)

	var overdue int64
	configs.DB.Model(&models.Invoice{}).
		Where("status = 'overdue' OR (due_date < ? AND status = 'pending')", todayStart).
		Count(&overdue)

	stats := map[string]interface{}{
		"TotalCustomers":     totalCustomers,
		"TotalInvoicesMonth": totalInvoicesMonth, 
		"DueToday":           dueToday,
		"Overdue":            overdue,
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
