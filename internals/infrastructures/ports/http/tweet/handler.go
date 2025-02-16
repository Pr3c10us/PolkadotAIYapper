package tweet

import (
	"github.com/Pr3c10us/boilerplate/internals/services/tweet"

	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/response"
	"github.com/gin-gonic/gin"
)

type Provider struct {
	Provider string `uri:"provider"  binding:"required"`
}

type Handler struct {
	services             tweet.Services
	environmentVariables *configs.EnvironmentVariables
}

func NewTweetHandler(service tweet.Services, environmentVariables *configs.EnvironmentVariables) Handler {
	return Handler{
		services:             service,
		environmentVariables: environmentVariables,
	}
}

func (handler *Handler) Topic(context *gin.Context) {
	_, err := handler.services.Tweet.Handle()
	if err != nil {
		//_ = context.Error(err)
		//fmt.Println(err)
		return
	}

	response.NewSuccessResponse("", nil, nil).Send(context)
}
