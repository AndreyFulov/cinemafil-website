package database

import (
	"context"
	"fmt"
	"log"

	"cinemafil-api/internal/config"
	"cinemafil-api/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	// Включаем расширение pg_trgm для быстрого нечеткого поиска по названию
	db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm;")

	// Автоматические миграции таблиц
	err = db.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Review{},
	)
	if err != nil {
		log.Fatalf("Ошибка выполнения AutoMigrate: %v", err)
	}

	log.Println("Подключение к PostgreSQL успешно установлено, миграции применены")
	return db
}

func ConnectRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Printf("Предупреждение: Redis недоступен (%v)", err)
	} else {
		log.Println("Подключение к Redis успешно установлено")
	}

	return client
}
