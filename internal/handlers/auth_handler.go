package handlers

import (
	"net/http"
	"reminder-tagihan/internal/configs"
	"reminder-tagihan/internal/models"
	"reminder-tagihan/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Error": c.Query("error"),
	})
}

func PerformLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if err := configs.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.Redirect(http.StatusFound, "/login?error=Invalid username or password")
		return
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		c.Redirect(http.StatusFound, "/login?error=Invalid username or password")
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	session.Save()

	// Log activity
	activity := models.ActivityLog{
		UserID:      user.ID,
		Action:      "login",
		Description: "User logged in",
	}
	configs.DB.Create(&activity)

	c.Redirect(http.StatusFound, "/dashboard")
}

func PerformLogout(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID != nil {
		uid := userID.(uint)
		// Log activity
		activity := models.ActivityLog{
			UserID:      uid,
			Action:      "logout",
			Description: "User logged out",
		}
		configs.DB.Create(&activity)
	}

	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}
