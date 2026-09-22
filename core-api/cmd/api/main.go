package main

import (
	"fmt"
	"log"

	"cinemafil-api/internal/config"
	"cinemafil-api/internal/database"
	"cinemafil-api/internal/handler"
	"cinemafil-api/internal/middleware"
	"cinemafil-api/internal/repository"
	"cinemafil-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("OMDb API Key загружен: '%s'", cfg.OMDbAPIKey)

	db := database.ConnectPostgres(cfg)
	rdb := database.ConnectRedis(cfg)

	// Инициализация репозиториев
	movieRepo := repository.NewMovieRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервисов
	movieService := service.NewMovieService(movieRepo, rdb, cfg.OMDbAPIKey)
	authService := service.NewAuthService(userRepo, cfg)

	// Инициализация хэндлеров
	movieHandler := handler.NewMovieHandler(movieService)
	authHandler := handler.NewAuthHandler(authService)

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

	api := app.Group("/api/v1")

	// Публичные маршруты аутентификации
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	// Фильмы
	movies := api.Group("/movies")
	movies.Get("/search", movieHandler.Search)
	movies.Get("/:slug", movieHandler.GetBySlug)
	movies.Get("/by-imdb/:imdb_id", movieHandler.GetByIMDbID)

	// Защищенные маршруты (требуют валидный Access Token)
	protected := api.Group("", middleware.Protected(cfg.JWTSecret))
	protected.Get("/users/me", authHandler.Me)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Сервер запущен на %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Ошибка запуска: %v", err)
	}
}
