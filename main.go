package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v3"
	serverv2 "github.com/jbarzegar/oci/internal/servers/v2"
)

func main() {
	engine := html.New("./views", ".html.tmpl")
	config := fiber.Config{
		Views: engine,
	}
	app := fiber.New(config)

	app.Get("/", func(c fiber.Ctx) error {
		allRoutes := c.App().GetRoutes(true)
		// Map routes
		return c.Render("index", fiber.Map{
			"Routes": allRoutes,
		})
	})

	app.Use(serverv2.New())

	app.Listen(":5000")
}
