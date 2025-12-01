package routes

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	formatted := time.Now().Format(time.RFC3339)
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("health func called at %s", formatted),
	})
	return
}

func healthRoute(r *gin.Engine) {
	healthRT := r.Group("/health")
	{
		healthRT.GET("/", health)
	}
}
