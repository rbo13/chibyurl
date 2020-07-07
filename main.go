package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
	"github.com/gofiber/limiter"
	"github.com/gofiber/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Environment File cannot be loaded: %v", err)
		return
	}

	var PORT = os.Getenv("PORT")

	server := fiber.New()

	// middlewares
	server.Use(
		middleware.Recover(),
		logger.New(),
		helmet.New(),
		// 5 requests per 15 seconds
		limiter.New(limiter.Config{
			Timeout: 15,
			Max:     5,
		}),
	)

	// serve static files
	server.Static("/", "./public")

	server.Listen(PORT)
}
