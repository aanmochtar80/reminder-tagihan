package main

import (
	"log"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/services"
)

func main() {
	configs.ConnectDB()
	services.InitWhatsApp()
	
	// Wait a bit to ensure WhatsApp connects (if it can)
	// We just want to see if ProcessReminders panics
	log.Println("Calling ProcessReminders...")
	services.ProcessReminders()
	log.Println("Done ProcessReminders!")
}
