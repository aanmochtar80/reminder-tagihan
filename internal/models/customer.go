package models

import "time"

type Customer struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Phone       string    `gorm:"not null" json:"phone"`
	Address     string    `json:"address"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Invoices    []Invoice `gorm:"foreignKey:CustomerID" json:"invoices,omitempty"`
}
