package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"
	
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/handlers"
	"reminder-tagihan/internal/middlewares"
	"reminder-tagihan/internal/models"
	"reminder-tagihan/internal/services"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, continuing with system ENV variables")
	}

	// Connect to database
	configs.ConnectDB()

	// Run database migrations and seed
	if err := models.Migrate(configs.DB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	if err := models.Seed(configs.DB); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Init services
	services.InitWhatsApp()
	services.InitCron()

	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init router
	r := gin.Default()

	// Setup Session
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "super-secret-key"
	}
	store := cookie.NewStore([]byte(secret))
	r.Use(sessions.Sessions("reminder_session", store))

	// Serve static files
	r.Static("/static", "./web/static")

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/**/*.html")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dashboard")
	})

	// Guest routes (Login)
	guest := r.Group("/")
	guest.Use(middlewares.GuestOnly())
	{
		guest.GET("/login", handlers.ShowLogin)
	}
	
	// Open POST route for login
	r.POST("/auth/login", handlers.PerformLogin)
	r.POST("/auth/logout", handlers.PerformLogout)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/whatsapp/qr", handlers.GetWhatsAppQR)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middlewares.AuthRequired())
	{
		protected.GET("/dashboard", handlers.ShowDashboard)
		protected.GET("/customers", handlers.ListCustomers)
		protected.POST("/customers", handlers.CreateCustomer)
		protected.GET("/invoices", handlers.ListInvoices)
		protected.POST("/invoices/generate", handlers.GenerateInvoices)
		protected.POST("/invoices/send-reminders", handlers.TriggerReminders)
		protected.POST("/invoices/:id/pay", handlers.MarkInvoicePaid)
		protected.GET("/whatsapp", handlers.ShowWhatsAppPage)
		protected.POST("/whatsapp/disconnect", handlers.DisconnectWhatsApp)
		protected.GET("/settings", handlers.ShowSettings)
		protected.POST("/settings", handlers.UpdateSettings)
		protected.GET("/logs", handlers.ShowLogs)
	}

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
