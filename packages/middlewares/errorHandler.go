package middlewares

import (
	"errors"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/Pr3c10us/boilerplate/packages/logger"
	"github.com/Pr3c10us/boilerplate/packages/response"
	"github.com/Pr3c10us/boilerplate/packages/validator"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func ErrorHandlerMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			var (
				pqErr           *pq.Error
				customError     *appError.CustomError
				validationError *validator.ValidationError
			)
			logger.LogWithFields("error", "Error handler message", zap.Error(err))

			switch {
			case errors.As(err.Err, &pqErr):
				{
					log.Print(pqErr.Code.Name())
					switch pqErr.Code {
					case "23505":
						response.ErrorResponse{
							StatusCode:   http.StatusConflict,
							Message:      "unique key value violated",
							ErrorMessage: pqErr.Detail,
						}.Send(c)
						return
					case "22P02":
						response.ErrorResponse{
							StatusCode:   http.StatusBadRequest,
							Message:      "invalid argument syntax",
							ErrorMessage: pqErr.Message,
						}.Send(c)
						return
					case "23503":
						response.ErrorResponse{
							StatusCode:   http.StatusBadRequest,
							Message:      "invalid foreign key identifier",
							ErrorMessage: pqErr.Detail,
						}.Send(c)
						return
					default:
						response.NewErrorResponse(pqErr).Send(c)
						return
					}
				}
			case errors.As(err.Err, &customError):
				response.ErrorResponse{
					StatusCode:   customError.StatusCode,
					Message:      customError.Message,
					ErrorMessage: customError.ErrorMessage,
				}.Send(c)
				return
			case errors.As(err.Err, &validationError):
				response.ErrorResponse{
					StatusCode:   validationError.StatusCode,
					Message:      validationError.Message,
					ErrorMessage: validationError.ErrorMessage,
				}.Send(c)
			default:
				response.NewErrorResponse(err).Send(c)
				return
			}
		}
	}
}
