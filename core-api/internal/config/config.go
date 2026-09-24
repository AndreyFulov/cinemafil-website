package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DBHost             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBPort             string
	DBSSLMode          string
	RedisAddr          string
	RedisPass          string
	JWTSecret          string
	FrontendURL        string
	APIBaseURL         string
	YandexClientID     string
	YandexClientSecret string
	OMDbAPIKey         string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
	}
	if err != nil {
		err = godotenv.Load("../../.env")
	}

	if err != nil {
		log.Println("Предупреждение: .env файл не найден, используются переменные окружения ОС")
	} else {
		log.Println("Успешно: файл .env загружен")
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "film_network"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:          getEnv("REDIS_PASSWORD", ""),
		JWTSecret:          getEnv("JWT_SECRET", "default_secret"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		APIBaseURL:         getEnv("API_BASE_URL", "http://localhost:8080"),
		YandexClientID:     getEnv("YANDEX_CLIENT_ID", ""),
		YandexClientSecret: getEnv("YANDEX_CLIENT_SECRET", ""),
		OMDbAPIKey:         getEnv("OMDB_API_KEY", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
