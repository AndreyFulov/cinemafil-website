package main

import (
	"fmt"
	"log"

	"cinemafil-api/internal/config"
	"cinemafil-api/internal/database"
	"cinemafil-api/internal/handler"
	"cinemafil-api/internal/repository"
	"cinemafil-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()

	db := database.ConnectPostgres(cfg)
	rdb := database.ConnectRedis(cfg)

	// Инициализация слоев
	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo, rdb)
	movieHandler := handler.NewMovieHandler(movieService)

	app := fiber.New(fiber.Config{
		AppName: "Film Network API v1.0",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Роуты API
	api := app.Group("/api/v1")

	movies := api.Group("/movies")
	movies.Get("/search", movieHandler.Search)
	movies.Get("/:slug", movieHandler.GetBySlug)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Сервер запущен на %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Ошибка запуска: %v", err)
	}
}
