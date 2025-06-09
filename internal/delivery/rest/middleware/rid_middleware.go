package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const headerXRequestID = "X-Request-ID"

const requestIDKey = "request_id"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerXRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set(requestIDKey, rid)
		c.Header(headerXRequestID, rid)
		c.Next()
	}
}

func GetRequestID(ctx context.Context) string {
	if requestId, ok := ctx.Value(requestIDKey).(string); ok {
		return requestId
	}
	return ""
}
