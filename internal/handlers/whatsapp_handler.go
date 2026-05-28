package handlers

import (
	"context"
	"net/http"
	"reminder-tagihan/internal/services"

	"github.com/gin-gonic/gin"
)

func ShowWhatsAppPage(c *gin.Context) {
	status := "Disconnected"
	if services.WAClient != nil {
		if services.WAClient.IsConnected() {
			if services.WAClient.IsLoggedIn() {
				status = "Connected & Logged In"
			} else {
				status = "Waiting for QR Scan"
			}
		}
	}

	c.HTML(http.StatusOK, "whatsapp.html", gin.H{
		"Status": status,
	})
}

func GetWhatsAppQR(c *gin.Context) {
	if services.WAClient != nil && services.WAClient.IsLoggedIn() {
		c.JSON(http.StatusOK, gin.H{"qr": "", "status": "logged_in"})
		return
	}

	// Try to get QR code string from channel
	select {
	case qrStr := <-services.QRChan:
		c.JSON(http.StatusOK, gin.H{"qr": qrStr, "status": "waiting"})
	default:
		c.JSON(http.StatusOK, gin.H{"qr": "", "status": "loading"})
	}
}

func DisconnectWhatsApp(c *gin.Context) {
	if services.WAClient != nil {
		services.WAClient.Logout(context.Background())
	}
	c.Redirect(http.StatusFound, "/whatsapp")
}
