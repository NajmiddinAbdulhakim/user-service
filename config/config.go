package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment       string // develop, staging, production
	PostgresHost      string
	PostgresPort      int
	PostgresDatabase  string
	PostgresUser      string
	PostgresPassword  string
	LogLevel          string
	RPCPort           string
	ReviewServiceHost string
	ReviewServicePort int
	PostServicePort int
	PostServiceHost string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(look("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(look("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(look("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(look("POSTGRES_DATABASE", "userdb"))
	c.PostgresUser = cast.ToString(look("POSTGRES_USER", "najmiddin"))
	c.PostgresPassword = cast.ToString(look("POSTGRES_PASSWORD", "1234"))

	c.PostServiceHost = cast.ToString(look("POST_SERVICE_HOST", "localhost"))
	c.PostServicePort = cast.ToInt(look("POST_SERVICE_PORT", 2222))
	c.LogLevel = cast.ToString(look("LOG_LEVEL", "debug"))

	c.RPCPort = cast.ToString(look("RPC_PORT", ":9999"))

	return c
}

func look(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
