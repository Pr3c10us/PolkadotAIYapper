package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type PostgresDB struct {
	Username string
	Password string
	Host     string
	Port     int
	Name     string
	SSLMode  string
}

type RedisCache struct {
	Address  string
	Password string
	Username string
}

type RedisKeys struct {
	VerificationCodeKey string
}

type AWSKeys struct {
	AccessID    string
	SecretKey   string
	Region      string
	Endpoint    string
	AWSFromMail string
}

type GoogleOAuth struct {
	GoogleKey      string
	GoogleSecret   string
	GoogleCallback string
}

type GithubOAuth struct {
	GithubKey      string
	GithubSecret   string
	GithubCallback string
}

type OAuthProvider struct {
	Google *GoogleOAuth
	Github *GithubOAuth
}

type SMTP struct {
	FromAddress string
	Host        string
	Port        int
	Username    string
	Password    string
}

type XDotCom struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessKey      string
	AccessSecret   string
	BearerToken    string
}

type EnvironmentVariables struct {
	Port                  string
	JWTSecret             string
	JWTMaxAge             time.Duration
	RefreshJWTSecret      string
	RefreshJWTMaxAge      time.Duration
	CookieSecret          string
	SessionSecret         string
	SessionMaxAge         int
	ProductionEnvironment bool
	AuthRedirectUrl       string
	ClientDomain          string
	ProjectName           string
	PostgresDB            *PostgresDB
	RedisCache            *RedisCache
	RedisKeys             *RedisKeys
	AWSKeys               *AWSKeys
	OAuthProvider         *OAuthProvider
	SMTP                  *SMTP
	OpenAIApiKey          string
	XDotCom               *XDotCom
}

func loadEnv() {
	rootPath := GetRootPath()
	err := godotenv.Load(rootPath + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func LoadEnvironment() *EnvironmentVariables {
	loadEnv()
	return &EnvironmentVariables{
		Port:                  getEnv("PORT", ":5000"),
		JWTSecret:             getEnvOrError("JWT_SECRET"),
		JWTMaxAge:             time.Second * time.Duration(getEnvAsInt("JWT_MAX_AGE", 60*15)),
		RefreshJWTSecret:      getEnvOrError("REFRESH_JWT_SECRET"),
		RefreshJWTMaxAge:      time.Second * time.Duration(getEnvAsInt("REFRESH_JWT_MAX_AGE", 60*60*24*31)),
		CookieSecret:          getEnvOrError("COOKIE_SECRET"),
		SessionSecret:         getEnvOrError("SESSIONS_SECRET"),
		SessionMaxAge:         getEnvAsInt("SESSION_MAX_AGE", 86400*300),
		ProductionEnvironment: getEnvAsBool("PRODUCTION_ENVIRONMENT", false),
		ClientDomain:          getEnv("CLIENT_DOMAIN", "localhost"),
		ProjectName:           getEnv("PROJECT_NAME", "rider"),
		PostgresDB: &PostgresDB{
			Username: getEnv("PG_DB_USERNAME", "postgres"),
			Password: getEnvOrError("PG_DB_PASSWORD"),
			Host:     getEnv("PG_DB_HOST", "127.0.0.1"),
			Port:     getEnvAsInt("PG_DB_PORT", 5432),
			Name:     getEnvOrError("PG_DB_NAME"),
			SSLMode:  getEnv("PG_SSL_MODE", "disable"),
		},
		RedisCache: &RedisCache{
			Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", "1234"),
			Username: getEnv("REDIS_USERNAME", "default"),
		},
		RedisKeys: &RedisKeys{
			VerificationCodeKey: getEnv("REDIS_VERIFICATION_CODE_KEY", "verification_code"),
		},
		OAuthProvider: &OAuthProvider{
			Google: &GoogleOAuth{
				GoogleKey:      getEnvOrError("GOOGLE_CLIENT_ID"),
				GoogleSecret:   getEnvOrError("GOOGLE_CLIENT_SECRET"),
				GoogleCallback: getEnvOrError("GOOGLE_CALLBACK_URL"),
			},
			Github: &GithubOAuth{
				GithubKey:      getEnvOrError("GITHUB_CLIENT_ID"),
				GithubSecret:   getEnvOrError("GITHUB_CLIENT_SECRET"),
				GithubCallback: getEnvOrError("GITHUB_CALLBACK_URL"),
			},
		},
		SMTP: &SMTP{
			FromAddress: getEnvOrError("SMTP_FROM_ADDRESS"),
			Host:        getEnvOrError("SMTP_HOST"),
			Port:        getEnvAsInt("SMTP_PORT", 587),
			Username:    getEnvOrError("SMTP_USERNAME"),
			Password:    getEnvOrError("SMTP_PASSWORD"),
		},
		OpenAIApiKey: getEnvOrError("OPENAI_API_KEY"),
		XDotCom: &XDotCom{
			ConsumerKey:    getEnvOrError("CONSUMERKEY"),
			ConsumerSecret: getEnvOrError("CONSUMERSECRET"),
			AccessKey:      getEnvOrError("ACCESS_KEY"),
			AccessSecret:   getEnvOrError("ACCESS_SECRET"),
			BearerToken:    getEnvOrError("BEARERTOKEN"),
		},
	}
}

func getEnvOrError(key string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	panic("Environment variable " + key + " not set")
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value, exist := os.LookupEnv(key)
	if exist {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			log.Panicf("Environment variable \"%v\" not set properly", key)
		}
		return valueInt
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	value, exist := os.LookupEnv(key)
	if exist {
		valueBool, err := strconv.ParseBool(value)
		if err != nil {
			log.Panicf("Environment variable \"%v\" not set properly", key)
		}
		return valueBool
	}
	return fallback
}
