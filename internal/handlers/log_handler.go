package handlers

import (
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"

	"github.com/gin-gonic/gin"
)

func ShowLogs(c *gin.Context) {
	var waLogs []models.WhatsAppLog
	configs.DB.Preload("Customer").Order("created_at desc").Limit(50).Find(&waLogs)
	
	var actLogs []models.ActivityLog
	configs.DB.Preload("User").Order("created_at desc").Limit(50).Find(&actLogs)

	c.HTML(http.StatusOK, "logs.html", gin.H{
		"WALogs":  waLogs,
		"ActLogs": actLogs,
	})
}
