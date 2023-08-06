package config

import (
	"os"
	"time"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	Server      struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	CDN struct {
		AwsAccessKeyID     string
		AwsSecretAccessKey string
		AwsEndpoint        string
		BucketName         string
		CdnBaseUrl         string
	}
	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SSLMode  string
	}
	Token struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
	}
	Context struct {
		Timeout string
	}
}

func NewConfig() (*Config, error) {

	config := Config{}
	// general initialization
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	// server initialization
	config.Server.Host = getEnv("SERVER_HOST", "")
	config.Server.Port = getEnv("SERVER_PORT", ":8080")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// cdn init
	config.CDN.AwsAccessKeyID = getEnv("AWS_ACCESS_KEY_ID", "")
	config.CDN.AwsSecretAccessKey = getEnv("AWS_SECRET_ACCESS_KEY", "")
	config.CDN.AwsEndpoint = getEnv("AWS_END_POINT", "")
	config.CDN.BucketName = getEnv("BUCKET_NAME", "")
	config.CDN.CdnBaseUrl = getEnv("CDN_BASE_URL", "")

	// db initialization
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "abclinic")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "postgres")
	config.DB.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")

	// token init
	// token initialization
	config.Token.Secret = getEnv("TOKEN_SECRET", "haCBL1+6BsR6f0ZyF6AOWh0mqouoqkl9KpNK9YhuVlA=")

	// access ttl parse
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "1h"))
	if err != nil {
		return nil, err
	}
	// refresh ttl parse
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "24h"))
	if err != nil {
		return nil, err
	}
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
