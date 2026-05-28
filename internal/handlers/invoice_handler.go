package handlers

import (
	"fmt"
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"
	"reminder-tagihan/internal/services"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ListInvoices(c *gin.Context) {
	var invoices []models.Invoice
	
	statusFilter := c.Query("status")
	
	query := configs.DB.Preload("Customer").Order("created_at desc")
	if statusFilter != "" && statusFilter != "all" {
		query = query.Where("status = ?", statusFilter)
	}
	
	query.Find(&invoices)

	session := sessions.Default(c)
	flashMsg := session.Get("flash")
	if flashMsg != nil {
		session.Delete("flash")
		session.Save()
	}

	c.HTML(http.StatusOK, "invoices.html", gin.H{
		"Invoices":     invoices,
		"StatusFilter": statusFilter,
		"FlashMessage": flashMsg,
	})
}

func MarkInvoicePaid(c *gin.Context) {
	id := c.Param("id")
	
	var invoice models.Invoice
	if err := configs.DB.First(&invoice, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/invoices")
		return
	}

	invoice.Status = "paid"
	configs.DB.Save(&invoice)

	// Create payment record automatically
	payment := models.Payment{
		InvoiceID:   invoice.ID,
		Amount:      invoice.Amount,
		PaymentDate: time.Now(),
		Method:      "Manual Marking",
	}
	configs.DB.Create(&payment)

	session := sessions.Default(c)
	session.Set("flash", fmt.Sprintf("Tagihan %s berhasil ditandai Lunas.", invoice.InvoiceNumber))
	session.Save()

	// Log activity
	activity := models.ActivityLog{
		UserID:      session.Get("user_id").(uint),
		Action:      "pay_invoice",
		Description: fmt.Sprintf("Marked invoice %s as paid", invoice.InvoiceNumber),
	}
	configs.DB.Create(&activity)

	c.Redirect(http.StatusFound, "/invoices")
}

func GenerateInvoices(c *gin.Context) {
	services.GenerateMonthlyInvoices()
	c.Redirect(http.StatusFound, "/invoices")
}

func TriggerReminders(c *gin.Context) {
	services.ProcessReminders()
	c.Redirect(http.StatusFound, "/invoices")
}
