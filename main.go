package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
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
	)

	// serve static files
	server.Static("/", "./public")

	server.Listen(PORT)
}
