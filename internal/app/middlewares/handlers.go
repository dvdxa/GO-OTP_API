package middlewares

import "github.com/gin-gonic/gin"

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "method not allowed or missing handler function"})
	}
}
