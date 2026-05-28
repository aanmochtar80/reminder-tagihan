package models

import "time"

type Invoice struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	InvoiceNumber string    `gorm:"uniqueIndex;not null" json:"invoice_number"`
	CustomerID    uint      `gorm:"index;not null" json:"customer_id"`
	Customer      Customer  `gorm:"foreignKey:CustomerID" json:"customer"`
	Amount        float64   `gorm:"not null" json:"amount"`
	IssueDate     time.Time `json:"issue_date"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `gorm:"default:'pending';comment:'pending, paid, overdue'" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Relations
	Payments      []Payment `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
}
