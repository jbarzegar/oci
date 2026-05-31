package main

import "github.com/gofiber/fiber/v3"

func main() {
	config := fiber.Config{}
	app := fiber.New(config)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello there")
	})

	app.Listen(":5000")
}
