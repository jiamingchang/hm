package pkg

import (
	"github.com/gin-gonic/gin"
	"hm/models"
)

func Current(c *gin.Context) (user *models.User){
	return c.MustGet("user").(*models.User)
}
