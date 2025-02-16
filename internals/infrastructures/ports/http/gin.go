package http

import (
	authentication2 "github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/internals/infrastructures/ports/http/authentication"
	"github.com/Pr3c10us/boilerplate/internals/infrastructures/ports/http/tweet"
	"github.com/Pr3c10us/boilerplate/internals/services"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/logger"
	"github.com/Pr3c10us/boilerplate/packages/middlewares"
	"github.com/Pr3c10us/boilerplate/packages/response"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type GinServer struct {
	Services    *services.Services
	Engine      *gin.Engine
	Logger      logger.Logger
	Environment *configs.EnvironmentVariables
}

func NewGinServer(services *services.Services, logger logger.Logger, environment *configs.EnvironmentVariables) *GinServer {
	ginServer := &GinServer{
		Services:    services,
		Engine:      gin.Default(),
		Logger:      logger,
		Environment: environment,
	}

	cookieStore := cookie.NewStore([]byte(environment.CookieSecret))
	gothic.Store = cookieStore

	// Middlewares
	ginServer.Engine.Use(middlewares.RequestLoggingMiddleware(logger))
	ginServer.Engine.Use(sessions.Sessions("auth", cookieStore))
	ginServer.Engine.Use(middlewares.ErrorHandlerMiddleware(logger))
	ginServer.Engine.NoRoute(middlewares.RouteNotFoundMiddleware())

	ginServer.Health()
	ginServer.SecureHealth()
	ginServer.Authentication()
	ginServer.Tweet()

	return ginServer
}
func (server *GinServer) SecureHealth() {
	server.Engine.GET("/secure/health", middlewares.UserAuthorizationMiddleware(server.Services.AuthenticationServices, server.Environment), func(c *gin.Context) {
		user := c.MustGet("user").(*authentication2.User)
		server.Logger.LogWithFields("info", user.FullName, user)
		response.NewSuccessResponse("server up!!!", nil, nil).Send(c)
	})
}

func (server *GinServer) Health() {
	server.Engine.GET("/health", func(c *gin.Context) {
		response.NewSuccessResponse("server up!!!", nil, nil).Send(c)
	})
}

func (server *GinServer) Authentication() {
	handler := authentication.NewAuthenticationHandler(server.Services.AuthenticationServices, server.Environment)
	oauthRoute := server.Engine.Group("/auth/:provider")
	{
		oauthRoute.GET("/", handler.InitiateAuth)
		oauthRoute.GET("/callback", handler.Callback)
	}
	authRoute := server.Engine.Group("/api/v1/authentication")
	{
		authRoute.POST("/signup", handler.Registration)
		authRoute.POST("/login", handler.Login)
		authRoute.GET("/token", handler.GetAccessToken)
		authRoute.POST("/verify/resend", handler.ResendCode)
		authRoute.POST("/verify/email", handler.VerifyCode)
	}
}
func (server *GinServer) Tweet() {
	handler := tweet.NewTweetHandler(server.Services.TweetService, server.Environment)
	route := server.Engine.Group("/tweet/topic")
	{
		route.GET("/", handler.Topic)
	}
}

func (server *GinServer) Run() {
	err := server.Engine.Run(server.Environment.Port)
	if err != nil {
		server.Logger.Log("panic", "failed to start server")
	}
}
