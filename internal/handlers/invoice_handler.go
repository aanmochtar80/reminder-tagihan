package handlers

import (
	"fmt"
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"
	"reminder-tagihan/internal/services"
	"time"

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

	var activeCustomers []models.Customer
	configs.DB.Where("is_active = ?", true).Order("name asc").Find(&activeCustomers)

	session := sessions.Default(c)
	flashMsg := session.Get("flash")
	if flashMsg != nil {
		session.Delete("flash")
		session.Save()
	}

	c.HTML(http.StatusOK, "invoices.html", gin.H{
		"Invoices":     invoices,
		"Customers":    activeCustomers,
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

func CreateInvoice(c *gin.Context) {
	customerID, _ := strconv.Atoi(c.PostForm("customer_id"))
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	dueDate, _ := time.Parse("2006-01-02", c.PostForm("due_date"))
	
	invoice := models.Invoice{
		InvoiceNumber: fmt.Sprintf("INV-%s-%d", time.Now().Format("20060102150405"), customerID),
		CustomerID:    uint(customerID),
		ItemName:      c.PostForm("item_name"),
		Amount:        amount,
		IssueDate:     time.Now(),
		DueDate:       dueDate,
		IsRecurring:   c.PostForm("is_recurring") == "on",
		Status:        "pending",
	}

	configs.DB.Create(&invoice)
	
	session := sessions.Default(c)
	session.Set("flash", "Tagihan berhasil dibuat.")
	session.Save()
	
	c.Redirect(http.StatusFound, "/invoices")
}
