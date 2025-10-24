package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		gin.DefaultWriter.Write([]byte(
			fmt.Sprintf("Method: %s | Path: %s | Status: %d | Duration: %v\n",
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				duration,
			),
		))
	}
}
