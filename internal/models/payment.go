package models

import "time"

type Payment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	InvoiceID   uint      `gorm:"index;not null" json:"invoice_id"`
	Invoice     Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice"`
	Amount      float64   `gorm:"not null" json:"amount"`
	PaymentDate time.Time `gorm:"not null" json:"payment_date"`
	Method      string    `json:"method"` // e.g. Transfer, Cash
	ProofImage  string    `json:"proof_image"` // path to uploaded image
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}
