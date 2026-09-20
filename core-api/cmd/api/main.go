package main

import (
	"fmt"
	"log"

	"cinemafil-api/internal/config"
	"cinemafil-api/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Загрузка конфигурации
	cfg := config.LoadConfig()

	// 2. Инициализация хранилищ
	db := database.ConnectPostgres(cfg)
	_ = db // будет передан в сервисы/репозитории
	rdb := database.ConnectRedis(cfg)
	_ = rdb

	// 3. Создание инстанса Fiber
	app := fiber.New(fiber.Config{
		AppName: "Film Network API v1.0",
	})

	// 4. Глобальные Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // URL Next.js
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// 5. Тестовый эндпоинт проверки здоровья
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "Film Network API is running",
		})
	})

	// 6. Запуск сервера
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Сервер запущен на %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
