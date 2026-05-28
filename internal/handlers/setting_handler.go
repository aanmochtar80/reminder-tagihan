package handlers

import (
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"
	"reminder-tagihan/internal/utils"

	"github.com/gin-contrib/sessions"
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
		"FlashMsg": sessions.Default(c).Flashes(),
	})
	sessions.Default(c).Save()
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

func UpdatePassword(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	newPassword := c.PostForm("new_password")
	confirmPassword := c.PostForm("confirm_password")

	if newPassword != confirmPassword {
		session.AddFlash("Password baru dan konfirmasi password tidak cocok!")
		session.Save()
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		session.AddFlash("Gagal mengamankan password baru.")
		session.Save()
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	var user models.User
	if err := configs.DB.First(&user, userID).Error; err == nil {
		user.Password = hash
		configs.DB.Save(&user)
		session.AddFlash("Password admin berhasil diubah!")
	} else {
		session.AddFlash("Gagal menemukan user.")
	}

	session.Save()
	c.Redirect(http.StatusFound, "/settings")
}
