package main

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// initializeApp initializes the Fiber app
func initializeApp() *fiber.App {
	engine := html.New("./views", ".html")

	// Add unescape function
	engine.AddFunc("unescape", func(s string) template.HTML {
		return template.HTML(s)
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Define routes
	defineRoutes(app)

	return app
}
