package middleware

import (
	"time"
	"vrs-api/internal/util/logger"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLogger()
		startTime := time.Now()
		c.Next()
		elapseTime := time.Since(startTime)
		param := map[string]any{
			"status_code": c.Writer.Status(),
			"method":      c.Request.Method,
			"latency":     elapseTime,
			"path":        c.Request.URL,
		}

		log.Info("Request id: " + GetRequestID(c))
		log.WithFields(param).Info("Request information:")
	}
}
