package services

import (
	"bytes"
	"fmt"
	"log"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"
	"strconv"
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

func formatCurrency(amount float64) string {
	s := strconv.FormatFloat(amount, 'f', 0, 64)
	n := len(s)
	if n <= 3 {
		return s
	}
	out := make([]byte, 0, n+(n-1)/3)
	for i := 0; i < n; i++ {
		out = append(out, s[i])
		if (n-i-1)%3 == 0 && i != n-1 {
			out = append(out, '.')
		}
	}
	return string(out)
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
	tmplStr := "Halo *{{.nama}}* 👋,\n\nTagihan layanan *{{.layanan}}* Anda sebesar *Rp {{.nominal}}* 💰 akan jatuh tempo pada *{{.tanggal}}* 📅.\n\nMohon segera melakukan pembayaran. _(Abaikan pesan ini jika sudah membayar)._\n\n💳 *Pembayaran bisa melalui:*\n🏦 Bank BNI: 0456659645\n🏦 Bank BCA: 7974208688\n🏦 Bank Jago: 103633049732\n📱 DANA: 085255627216\n📱 OVO: 085255627216\n👤 *A.N Hamzah Mochtar*\n\nTerima kasih 🙏"
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
				"layanan": inv.ItemName,
				"nominal": formatCurrency(inv.Amount),
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

func ProcessSelectedReminders(invoiceIDs []string) {
	if WAClient == nil || !WAClient.IsConnected() || !WAClient.IsLoggedIn() {
		log.Println("Manual reminder skipped: WhatsApp not connected")
		return
	}

	var invoices []models.Invoice
	configs.DB.Preload("Customer").Where("id IN ?", invoiceIDs).Find(&invoices)

	var setting models.Setting
	tmplStr := "Halo *{{.nama}}* 👋,\n\nTagihan layanan *{{.layanan}}* Anda sebesar *Rp {{.nominal}}* 💰 akan jatuh tempo pada *{{.tanggal}}* 📅.\n\nMohon segera melakukan pembayaran. _(Abaikan pesan ini jika sudah membayar)._\n\n💳 *Pembayaran bisa melalui:*\n🏦 Bank BNI: 0456659645\n🏦 Bank BCA: 7974208688\n🏦 Bank Jago: 103633049732\n📱 DANA: 085255627216\n📱 OVO: 085255627216\n👤 *A.N Hamzah Mochtar*\n\nTerima kasih 🙏"
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

		var buf bytes.Buffer
		data := map[string]interface{}{
			"nama":    inv.Customer.Name,
			"layanan": inv.ItemName,
			"nominal": formatCurrency(inv.Amount),
			"tanggal": inv.DueDate.Format("02 Jan 2006"),
			"invoice": inv.InvoiceNumber,
		}
		tmpl.Execute(&buf, data)

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

		logEntry := models.WhatsAppLog{
			CustomerID:   inv.CustomerID,
			InvoiceID:    &inv.ID,
			Type:         "manual_reminder",
			Message:      buf.String(),
			Status:       status,
			ErrorMessage: errMsg,
		}
		configs.DB.Create(&logEntry)
	}
}

// Helper to generate monthly invoices (Can be run on the 1st of every month)
func GenerateMonthlyInvoices() {
	var recurringInvoices []models.Invoice
	configs.DB.Where("is_recurring = ?", true).Order("due_date desc").Find(&recurringInvoices)

	now := time.Now()
	processed := make(map[string]bool)

	for _, inv := range recurringInvoices {
		key := fmt.Sprintf("%d_%s", inv.CustomerID, inv.ItemName)
		if processed[key] {
			continue
		}
		processed[key] = true

		dueDate := time.Date(now.Year(), now.Month(), inv.DueDate.Day(), 0, 0, 0, 0, now.Location())

		var count int64
		configs.DB.Model(&models.Invoice{}).
			Where("customer_id = ? AND item_name = ? AND strftime('%Y-%m', due_date) = ?", 
                  inv.CustomerID, inv.ItemName, dueDate.Format("2006-01")).
			Count(&count)

		if count == 0 {
			newInv := models.Invoice{
				InvoiceNumber: fmt.Sprintf("INV-%s-%d-%d", dueDate.Format("200601"), inv.CustomerID, time.Now().UnixNano()%10000),
				CustomerID:    inv.CustomerID,
				ItemName:      inv.ItemName,
				Amount:        inv.Amount,
				IsRecurring:   true,
				IssueDate:     now,
				DueDate:       dueDate,
				Status:        "pending",
			}
			configs.DB.Create(&newInv)
		}
	}
}
