package authentication

import (
	"errors"
	"fmt"

	authentication2 "github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/internals/services/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/middlewares"
	"github.com/Pr3c10us/boilerplate/packages/response"
	"github.com/Pr3c10us/boilerplate/packages/utils"
	"github.com/Pr3c10us/boilerplate/packages/validator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

type Provider struct {
	Provider string `uri:"provider"  binding:"required"`
}

type Handler struct {
	services             authentication.Services
	environmentVariables *configs.EnvironmentVariables
}

func NewAuthenticationHandler(service authentication.Services, environmentVariables *configs.EnvironmentVariables) Handler {
	return Handler{
		services:             service,
		environmentVariables: environmentVariables,
	}
}

func (handler *Handler) Registration(context *gin.Context) {
	var addUserRequest authentication2.AddUserParams
	if err := context.ShouldBind(&addUserRequest); err != nil {
		err = validator.ValidateRequest(err)
		fmt.Println(err)
		_ = context.Error(err)
		return
	}
	err := handler.services.CreateUser.Handle(&addUserRequest)
	if err != nil {
		_ = context.Error(err)
		fmt.Println(err)
		return
	}

	response.NewSuccessResponse("account created successfully, check email for verification code", nil, nil).Send(context)
}

func (handler *Handler) Login(context *gin.Context) {
	var loginParams authentication2.LoginParams
	if err := context.ShouldBind(&loginParams); err != nil {
		err = validator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}

	user, err := handler.services.Login.Handle(&loginParams)
	if err != nil {
		_ = context.Error(err)
		return
	}

	var token, refreshToken string
	token, err = utils.CreateUserToken(user, handler.environmentVariables.JWTSecret, handler.environmentVariables.JWTMaxAge)
	if err != nil {
		_ = context.Error(err)
		return
	}
	refreshToken, err = utils.CreateUserRefreshToken(user, handler.environmentVariables.RefreshJWTSecret, handler.environmentVariables.RefreshJWTMaxAge)
	if err != nil {
		_ = context.Error(err)
		return
	}

	session := sessions.Default(context)
	session.Set("token", token)
	session.Set("refreshToken", refreshToken)
	session.Options(sessions.Options{MaxAge: handler.environmentVariables.SessionMaxAge, Path: "/"})
	if err = session.Save(); err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"user": user, "token": token, "refreshToken": refreshToken}, nil).Send(context)
}

func (handler *Handler) VerifyCode(context *gin.Context) {
	var verifyCodeParams authentication2.VerifyCodeParams
	if err := context.ShouldBind(&verifyCodeParams); err != nil {
		err = validator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}

	err := handler.services.VerifyCode.Handle(&verifyCodeParams)
	if err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse("email verified", nil, nil).Send(context)
}

func (handler *Handler) ResendCode(context *gin.Context) {
	var resendCodeParams authentication2.ResendCodeParams
	if err := context.ShouldBind(&resendCodeParams); err != nil {
		err = validator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}

	err := handler.services.ResendCode.Handle(resendCodeParams.Email)
	if err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse("code send", nil, nil).Send(context)
}

func (handler *Handler) GetAccessToken(context *gin.Context) {
	var refreshToken string
	var err error

	// Get refresh token cookie or header if cookie fail
	refreshToken = middlewares.GetCookieKey(context, "refreshToken")
	if refreshToken == "" {
		refreshToken, err = middlewares.GetAuthorizationToken(context)
		if err != nil {
			_ = context.Error(err)
			return
		}
	}

	// Get token claim to retrieve token info
	var claims *utils.Claims
	claims, err = utils.DecryptUserToken(refreshToken, handler.environmentVariables.RefreshJWTSecret)
	if err != nil {
		_ = context.Error(err)
		return
	}

	// Fetch user information
	var user *authentication2.User
	var id uuid.UUID
	id, err = uuid.Parse(claims.ID)
	user, err = handler.services.GetUserDetails.Handle(&authentication2.GetUserParams{
		ID: id,
	})
	if err != nil {
		_ = context.Error(err)
		context.Abort()
		return
	}

	// If refresh token version is not the same as claim abort
	if user.RefreshTokenVersion != claims.Version {
		_ = context.Error(appError.Unauthorized(errors.New("token is expired")))
		return
	}

	// create new user token and refresh token
	var token string
	token, err = utils.CreateUserToken(user, handler.environmentVariables.JWTSecret, handler.environmentVariables.JWTMaxAge)
	if err != nil {
		_ = context.Error(err)
		return
	}
	refreshToken, err = utils.CreateUserRefreshToken(user, handler.environmentVariables.RefreshJWTSecret, handler.environmentVariables.RefreshJWTMaxAge)
	if err != nil {
		_ = context.Error(err)
		return
	}

	// store tokens into session
	session := sessions.Default(context)
	session.Set("token", token)
	session.Set("refreshToken", refreshToken)
	session.Options(sessions.Options{MaxAge: handler.environmentVariables.SessionMaxAge, Path: "/"})
	if err = session.Save(); err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"user": user, "token": token}, nil).Send(context)

}

func (handler *Handler) InitiateAuth(context *gin.Context) {
	var provider Provider
	if err := context.ShouldBindUri(&provider); err != nil {
		err = validator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}
	q := context.Request.URL.Query()
	q.Add("provider", provider.Provider)
	context.Request.URL.RawQuery = q.Encode()
	if gothUser, err := gothic.CompleteUserAuth(context.Writer, context.Request); err == nil {
		fmt.Println(gothUser, "-----------------------------------------------------")
		var user *authentication2.User
		user, err = handler.services.Authenticate.Handle(&gothUser)
		if err != nil {
			_ = context.Error(err)
			return
		}

		var token, refreshToken string
		token, err = utils.CreateUserToken(user, handler.environmentVariables.JWTSecret, handler.environmentVariables.JWTMaxAge)
		if err != nil {
			_ = context.Error(err)
			return
		}
		refreshToken, err = utils.CreateUserRefreshToken(user, handler.environmentVariables.RefreshJWTSecret, handler.environmentVariables.RefreshJWTMaxAge)
		if err != nil {
			_ = context.Error(err)
			return
		}

		session := sessions.Default(context)
		session.Set("token", token)
		session.Set("refreshToken", refreshToken)
		session.Options(sessions.Options{MaxAge: handler.environmentVariables.SessionMaxAge, Path: "/"})
		if err = session.Save(); err != nil {
			_ = context.Error(err)
			return
		}
		response.NewSuccessResponse("", gin.H{"user": user, "token": token, "refreshToken": refreshToken}, nil).Send(context)

	} else {
		fmt.Println(err, "------------------------------------")
		gothic.BeginAuthHandler(context.Writer, context.Request)
	}
}

func (handler *Handler) Callback(context *gin.Context) {
	var provider Provider
	if err := context.ShouldBindUri(&provider); err != nil {
		err = validator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}
	q := context.Request.URL.Query()
	q.Add("provider", provider.Provider)
	context.Request.URL.RawQuery = q.Encode()
	gothUser, err := gothic.CompleteUserAuth(context.Writer, context.Request)
	if err != nil {
		_ = context.Error(err)
		return
	}

	var user *authentication2.User
	user, err = handler.services.Authenticate.Handle(&gothUser)
	if err != nil {
		_ = context.Error(err)
		return
	}

	var token, refreshToken string
	token, err = utils.CreateUserToken(user, handler.environmentVariables.JWTSecret, handler.environmentVariables.JWTMaxAge)
	if err != nil {
		_ = context.Error(err)
		return
	}
	refreshToken, err = utils.CreateUserRefreshToken(user, handler.environmentVariables.RefreshJWTSecret, handler.environmentVariables.RefreshJWTMaxAge)
	if err != nil {
		_ = context.Error(err)
		return
	}

	session := sessions.Default(context)
	session.Set("token", token)
	session.Set("refreshToken", refreshToken)
	session.Options(sessions.Options{MaxAge: handler.environmentVariables.SessionMaxAge, Path: "/"})
	if err = session.Save(); err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse("", gin.H{"user": user, "token": token, "refreshToken": refreshToken}, nil).Send(context)
}
