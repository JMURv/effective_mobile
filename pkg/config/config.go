package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Server          *ServerConfig
	DB              *DBConfig
	ExternalAPIPort int
}

type ServerConfig struct {
	Mode   string
	Port   int
	Scheme string
	Domain string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("No .env file found")
	}

	return &Config{
		Server: &ServerConfig{
			Mode:   getEnv("SERVER_MODE", "dev"),
			Port:   getEnvAsInt("SERVER_PORT", 8080),
			Scheme: getEnv("SERVER_SCHEME", "http"),
			Domain: getEnv("SERVER_DOMAIN", "localhost"),
		},
		DB: &DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Database: getEnv("DB_NAME", "db"),
		},
		ExternalAPIPort: getEnvAsInt("EXTERNAL_API_PORT", 8081),
	}
}

func getEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(getEnv(key, "")); err == nil {
		return val
	}
	return defaultVal
}
