package models

import "time"

type Customer struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Phone       string    `gorm:"not null" json:"phone"`
	Address     string    `json:"address"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	BillAmount  float64   `gorm:"not null;default:0" json:"bill_amount"`
	DueDateDay  int       `gorm:"not null;default:1" json:"due_date_day"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Invoices    []Invoice `gorm:"foreignKey:CustomerID" json:"invoices,omitempty"`
}
