package models

import "time"

type ActivityLog struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"index"`
	User        User      `gorm:"foreignKey:UserID"`
	Action      string    `gorm:"not null"`
	Description string    
	CreatedAt   time.Time
}
