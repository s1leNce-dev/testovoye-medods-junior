package snippets

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleErrorJSONAnswer(c *gin.Context, statusCode int, errorMessage, errorLogMessage, logSource string) {
	log.Printf("[ERROR] %s %s\n", logSource, errorLogMessage)
	c.JSON(statusCode, gin.H{
		"error": errorMessage,
	})
}
