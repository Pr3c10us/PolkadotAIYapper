package middlewares

import (
	"github.com/Pr3c10us/boilerplate/packages/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func RequestLoggingMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		context.Next()

		logger.LogWithFields("info",
			"Request Information",
			zap.Int("status", context.Writer.Status()),
			zap.String("method", context.Request.Method),
			zap.String("path", context.Request.URL.Path),
			zap.Any("query", context.Request.URL.Query()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}
