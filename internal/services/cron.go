package services

import (
	"bytes"
	"fmt"
	"log"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"
	"strings"
	"text/template"
	"time"

	"github.com/robfig/cron/v3"
)

var ReminderCron *cron.Cron

func InitCron() {
	ReminderCron = cron.New()
	
	// Generate invoices every day at 01:00 AM (safe because it checks for duplicates)
	ReminderCron.AddFunc("0 1 * * *", GenerateMonthlyInvoices)
	
	// Check invoices every day at 08:00 AM
	ReminderCron.AddFunc("0 8 * * *", ProcessReminders)
	
	ReminderCron.Start()
	log.Println("Cron scheduler started")
}

func ProcessReminders() {
	if WAClient == nil || !WAClient.IsConnected() || !WAClient.IsLoggedIn() {
		log.Println("Cron skipped: WhatsApp not connected")
		return
	}

	now := time.Now()
	// Get all active customers with pending or overdue invoices
	var invoices []models.Invoice
	configs.DB.Preload("Customer").Where("status != ?", "paid").Find(&invoices)

	// Fetch custom template from Settings DB
	var setting models.Setting
	tmplStr := "Halo {{.nama}},\n\nTagihan layanan {{.layanan}} Anda sebesar Rp {{.nominal}} akan jatuh tempo pada {{.tanggal}}.\n\nMohon segera melakukan pembayaran. Abaikan pesan ini jika sudah membayar.\n\nTerima kasih."
	if err := configs.DB.Where("key = ?", "reminder_template").First(&setting).Error; err == nil && setting.Value != "" {
		tmplStr = setting.Value
	}
	tmpl, err := template.New("msg").Parse(tmplStr)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return
	}

	for _, inv := range invoices {
		if !inv.Customer.IsActive {
			continue
		}

		daysDiff := int(inv.DueDate.Sub(now).Hours() / 24)
		sendReminder := false
		reminderType := ""

		if daysDiff == 3 {
			sendReminder = true
			reminderType = "reminder_h3"
		} else if daysDiff == 0 {
			sendReminder = true
			reminderType = "reminder_hari_h"
		} else if daysDiff == -1 && inv.Status != "paid" {
			sendReminder = true
			reminderType = "overdue_h1"
			// Auto mark overdue
			inv.Status = "overdue"
			configs.DB.Save(&inv)
		}

		if sendReminder {
			var buf bytes.Buffer
			data := map[string]interface{}{
				"nama":    inv.Customer.Name,
				"layanan": inv.Customer.ServiceName,
				"nominal": inv.Amount,
				"tanggal": inv.DueDate.Format("02 Jan 2006"),
				"invoice": inv.InvoiceNumber,
			}
			tmpl.Execute(&buf, data)
			
			// Format phone number (replace leading 0 with 62)
			phone := inv.Customer.Phone
			if strings.HasPrefix(phone, "0") {
				phone = "62" + phone[1:]
			}
			jid := phone + "@s.whatsapp.net"
			
			err := SendMessage(jid, buf.String())
			
			status := "success"
			errMsg := ""
			if err != nil {
				status = "failed"
				errMsg = err.Error()
			}
			
			log := models.WhatsAppLog{
				CustomerID:   inv.CustomerID,
				InvoiceID:    &inv.ID,
				Type:         reminderType,
				Message:      buf.String(),
				Status:       status,
				ErrorMessage: errMsg,
			}
			configs.DB.Create(&log)
		}
	}
}

// Helper to generate monthly invoices (Can be run on the 1st of every month)
func GenerateMonthlyInvoices() {
	var customers []models.Customer
	configs.DB.Where("is_active = ?", true).Find(&customers)

	now := time.Now()
	for _, cust := range customers {
		// Calculate due date for current month
		dueDate := time.Date(now.Year(), now.Month(), cust.DueDateDay, 0, 0, 0, 0, now.Location())
		if dueDate.Before(now) {
			// If due date is already passed for this month, skip or create for next month
			// For simplicity, let's just generate for next month if the day has passed significantly
			// Let's keep it simple: always generate for the current month's due date
		}

		// Check if invoice for this customer in this month already exists
		var count int64
		configs.DB.Model(&models.Invoice{}).
			Where("customer_id = ? AND strftime('%Y-%m', due_date) = ?", cust.ID, dueDate.Format("2006-01")).
			Count(&count)

		if count == 0 {
			inv := models.Invoice{
				InvoiceNumber: fmt.Sprintf("INV-%s-%04d", dueDate.Format("200601"), cust.ID),
				CustomerID:    cust.ID,
				Amount:        cust.BillAmount,
				IssueDate:     now,
				DueDate:       dueDate,
				Status:        "pending",
			}
			configs.DB.Create(&inv)
		}
	}
}
