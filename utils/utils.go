package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// ReturnData :
func ReturnData(code int, msg string, status int, c *gin.Context) {
	c.JSON(code, gin.H{
		"status":  status,
		"message": msg,
	})
}

// ReturnJSON :
func ReturnJSON(val interface{}) string {
	b, _ := json.MarshalIndent(val, "", "    ")
	return string(b)
}
