package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

//go:embed frontend/dist/*
var embedDist embed.FS // output of vite

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDist),
		PathPrefix: "frontend/dist",
	}))
	app.Get("/api", func(c *fiber.Ctx) error {
		// panic("I'm an error")
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":9000"))
}