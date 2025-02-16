package adapters

import (
	"database/sql"
	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/internals/domains/cache"
	"github.com/Pr3c10us/boilerplate/internals/domains/email"
	"github.com/Pr3c10us/boilerplate/internals/domains/embedding"
	"github.com/Pr3c10us/boilerplate/internals/domains/llm"
	"github.com/Pr3c10us/boilerplate/internals/domains/sms"
	"github.com/Pr3c10us/boilerplate/internals/domains/xdotcom"
	authentication2 "github.com/Pr3c10us/boilerplate/internals/infrastructures/adapters/authentication"
	cache2 "github.com/Pr3c10us/boilerplate/internals/infrastructures/adapters/cache"
	email2 "github.com/Pr3c10us/boilerplate/internals/infrastructures/adapters/email"
	embedding2 "github.com/Pr3c10us/boilerplate/internals/infrastructures/adapters/embedding"
	openai2 "github.com/Pr3c10us/boilerplate/internals/infrastructures/adapters/llm/openai"
	xdotcom2 "github.com/Pr3c10us/boilerplate/internals/infrastructures/adapters/xdotcom"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/Pr3c10us/boilerplate/packages/logger"
	"github.com/openai/openai-go"
	"github.com/redis/go-redis/v9"
)

type AdapterDependencies struct {
	Logger               logger.Logger
	EnvironmentVariables *configs.EnvironmentVariables
	DB                   *sql.DB
	Redis                *redis.Client
	OpenAI               *openai.Client
}

type Adapters struct {
	Logger                   logger.Logger
	EnvironmentVariables     *configs.EnvironmentVariables
	AuthenticationRepository authentication.Repository
	EmailRepository          email.Repository
	SMSRepository            sms.Repository
	CacheRepository          cache.Repository
	OpenAiRepository         llm.Repository
	EmbeddingRepository      embedding.Repository
	XDotComRepository        xdotcom.Repository
}

func NewAdapters(dependencies AdapterDependencies) *Adapters {
	return &Adapters{
		Logger:                   dependencies.Logger,
		EnvironmentVariables:     dependencies.EnvironmentVariables,
		AuthenticationRepository: authentication2.NewAuthenticationRepositoryPG(dependencies.DB),
		EmailRepository:          email2.NewGoMailEmailRepository(dependencies.EnvironmentVariables),
		CacheRepository:          cache2.NewRedisRepository(dependencies.Redis, dependencies.EnvironmentVariables),
		OpenAiRepository:         openai2.NewOpenAIRepository(dependencies.OpenAI),
		EmbeddingRepository:      embedding2.NewEmbedding(dependencies.DB),
		XDotComRepository:        xdotcom2.NewXDotComRepository(dependencies.EnvironmentVariables),
	}
}
