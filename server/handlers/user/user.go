package user

import (
	"github.com/gin-gonic/gin"
)

func GetUUID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	c.JSON(200, gin.H{
		"message": "success",
		"data":    userID,
	})
}
