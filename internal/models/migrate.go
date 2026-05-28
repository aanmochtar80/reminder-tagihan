package models

import (
	"log"

	"gorm.io/gorm"
)

// Migrate runs auto-migration for all models
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	err := db.AutoMigrate(
		&User{},
		&Setting{},
		&ActivityLog{},
		&Customer{},
		&Invoice{},
		&Payment{},
		&WhatsAppLog{},
	)
	
	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}
	
	log.Println("Migration completed successfully.")
	return nil
}
