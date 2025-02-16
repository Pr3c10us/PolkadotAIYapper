package middlewares

import (
	"errors"
	"fmt"
	authentication2 "github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/internals/services/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

func UserAuthorizationMiddleware(service authentication.Services, environmentVariables *configs.EnvironmentVariables) gin.HandlerFunc {
	return func(context *gin.Context) {
		var token string
		var err error
		token = GetCookieKey(context, "token")
		if token == "" {
			token, err = GetAuthorizationToken(context)
			if err != nil {
				_ = context.Error(err)
				context.Abort()
				return
			}
		}

		var claims *utils.Claims
		claims, err = utils.DecryptUserToken(token, environmentVariables.JWTSecret)
		if err != nil {
			_ = context.Error(err)
			context.Abort()
			return
		}

		var user *authentication2.User
		var id uuid.UUID
		id, err = uuid.Parse(claims.ID)
		user, err = service.GetUserDetails.Handle(&authentication2.GetUserParams{
			ID: id,
		})
		if err != nil {
			_ = context.Error(err)
			context.Abort()
			return
		}

		context.Set("user", user)
		context.Next()
	}
}

func GetAuthorizationToken(context *gin.Context) (string, error) {
	bearerToken := context.Request.Header.Get("Authorization")
	if bearerToken == "" {
		err := appError.Unauthorized(errors.New("authorization token missing"))
		return "", err
	}

	tokenParts := strings.Split(bearerToken, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		err := appError.Unauthorized(errors.New("invalid authorization header format"))
		return "", err
	}

	return tokenParts[1], nil
}

func GetCookieKey(context *gin.Context, key string) string {
	session := sessions.Default(context)
	token := session.Get(key)
	fmt.Println(token)
	if token == nil {
		return ""
	} else {
		return token.(string)
	}
}
