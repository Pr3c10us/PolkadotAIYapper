package middlewares

import (
	"errors"
	"fmt"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/gin-gonic/gin"
)

func RouteNotFoundMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		var err = errors.New(fmt.Sprintf("Route '%v' not found", context.Request.URL.Path))
		context.Error(appError.NotFound(err))
	}
}
