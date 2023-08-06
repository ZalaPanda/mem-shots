package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ZalaPanda/mem-shots/processor"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

//go:embed frontend/dist/*
var embedDist embed.FS

// func DefaultErrorHandler(c *fiber.Ctx, err error) error {
// 	return fiber.ErrBadRequest // https://docs.gofiber.io/guide/error-handling
// }

func main() {
	bp := processor.New()

	log.SetLevel(log.LevelTrace)
	log.SetOutput(os.Stdout)

	config := fiber.Config{
		// ErrorHandler: DefaultErrorHandler,
	}
	app := fiber.New(config)

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDist),
		PathPrefix: "frontend/dist",
	}))
	app.Get("/api", func(c *fiber.Ctx) error {
		log.Info("API request")
		// panic("I'm an error")
		return c.Status(fiber.StatusOK).SendString("Hello World!")
	})
	app.Get("/api/start", func(c *fiber.Ctx) error {
		bp.Start()
		return c.Status(fiber.StatusOK).SendString(bp.GetState())
	})
	app.Get("/api/progress", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")
		c.Status(fiber.StatusOK)
		listener := bp.Listen()
		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()
			for progress := range listener {
				str := fmt.Sprintf("data: %s %d\n\n", "hello", progress)
				pw.Write([]byte(str))
			}
			pw.Close()
		}()
		return c.SendStream(pr)
	})
	app.Get("/api/stop", func(c *fiber.Ctx) error {
		bp.Stop()
		return c.Status(fiber.StatusOK).SendString(bp.GetState())
	})

	log.Fatal(app.Listen(":9000"))
}
