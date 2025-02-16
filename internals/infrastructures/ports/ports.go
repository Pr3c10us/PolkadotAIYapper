package ports

import (
	"github.com/Pr3c10us/boilerplate/internals/infrastructures/ports/http"
	"github.com/Pr3c10us/boilerplate/internals/services"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/logger"
)

type Ports struct {
	GinServer *http.GinServer
}

func NewPorts(services *services.Services, logger logger.Logger, environment *configs.EnvironmentVariables) *Ports {
	return &Ports{
		GinServer: http.NewGinServer(services, logger, environment),
	}
}
