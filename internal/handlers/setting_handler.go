package handlers

import (
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"

	"github.com/gin-gonic/gin"
)

func ShowSettings(c *gin.Context) {
	var settings []models.Setting
	configs.DB.Find(&settings)
	
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"Settings": settingsMap,
	})
}

func UpdateSettings(c *gin.Context) {
	c.Request.ParseForm()
	
	for key, values := range c.Request.PostForm {
		val := values[0]
		
		var setting models.Setting
		if err := configs.DB.Where("key = ?", key).First(&setting).Error; err != nil {
			setting = models.Setting{Key: key, Value: val}
			configs.DB.Create(&setting)
		} else {
			setting.Value = val
			configs.DB.Save(&setting)
		}
	}

	c.Redirect(http.StatusFound, "/settings")
}
