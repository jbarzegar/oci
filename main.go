package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/template/html/v3"
	serverv2 "github.com/jbarzegar/oci/internal/servers/v2"
)

func main() {
	engine := html.New("./views", ".gohtml")
	config := fiber.Config{
		Views: engine,
	}
	app := fiber.New(config)
	// Init recovery handler.
	// enable stack traces so we can more easily debug
	app.Use(recoverer.New(recoverer.Config{
		Next:              nil,
		PanicHandler:      recoverer.DefaultPanicHandler,
		StackTraceHandler: recoverer.ConfigDefault.StackTraceHandler,
		EnableStackTrace:  true,
	}))
	app.Use(logger.New())

	app.Get("/", func(c fiber.Ctx) error {
		allRoutes := c.App().GetRoutes(true)
		// Map routes
		return c.Render("index", fiber.Map{
			"Routes": allRoutes,
		})
	})

	s, err := serverv2.New()
	if err != nil {
		log.Fatal(err)
	}
	app.Use(s)

	app.Listen(":5000")
}
