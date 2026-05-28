package models

import "time"

type WhatsAppLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CustomerID   uint      `gorm:"index" json:"customer_id"`
	Customer     Customer  `gorm:"foreignKey:CustomerID" json:"customer"`
	InvoiceID    *uint     `gorm:"index" json:"invoice_id,omitempty"` // optional, for reminder broadcasts that are not tied to specific invoice sometimes
	Type         string    `json:"type"` // e.g. reminder, overdue, payment, broadcast
	Message      string    `json:"message"`
	Status       string    `gorm:"default:'pending';comment:'pending, success, failed'" json:"status"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}
