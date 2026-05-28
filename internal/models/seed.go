package models

import (
	"log"
	"reminder-tagihan/internal/utils"

	"gorm.io/gorm"
)

// Seed initial data if necessary
func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&User{}).Count(&count)
	
	if count == 0 {
		log.Println("No users found. Seeding default admin user...")
		hash, _ := utils.HashPassword("admin")
		admin := User{
			Username: "admin",
			Password: hash,
			Name:     "Administrator",
			Role:     "admin",
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		log.Println("Default admin user created: admin / admin")
	}
	return nil
}
