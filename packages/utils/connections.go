package utils

import (
	"database/sql"
	"fmt"
	"github.com/Pr3c10us/boilerplate/packages/configs"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	redis "github.com/redis/go-redis/v9"
)

func NewPGConnection(env *configs.EnvironmentVariables) *sql.DB {
	// PG_DB instantiation
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		env.PostgresDB.Host,
		env.PostgresDB.Port,
		env.PostgresDB.Username,
		env.PostgresDB.Password,
		env.PostgresDB.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic("failed to instantiation DB connection")
	}
	err = db.Ping()
	if err != nil {
		panic("no connection could be made because the target machine actively refused it")
	}
	return db
}

func NewRedisClient(env *configs.EnvironmentVariables) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     env.RedisCache.Address,
		Password: env.RedisCache.Password,
		Username: env.RedisCache.Username,
		DB:       0,
	})
	return redisClient
}

func NewOpenAIClient(env *configs.EnvironmentVariables) *openai.Client {
	return openai.NewClient(
		option.WithAPIKey(env.OpenAIApiKey),
	)
}
