package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/template/html/v3"
	serverv2 "github.com/jbarzegar/oci/internal/servers/v2"
)

func main() {
	// client, err := blobstorage.NewS3Storer("lol", true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := client.WriteImage(
	// 	context.Background(),
	// 	"hello", "123",
	// 	strings.NewReader("Hello there"),
	// ); err != nil {
	// 	log.Fatal(err)
	// }

	engine := html.New("./views", ".gohtml")
	config := fiber.Config{
		Views: engine,
	}
	app := fiber.New(config)
	app.Use(logger.New())

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
