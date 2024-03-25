package main

import "github.com/gofiber/fiber/v2"

// defineRoutes defines routes for the Fiber app
func defineRoutes(app *fiber.App) {
	// Define routes
	app.Get("/", handleIndex)
	app.Get("/order/:id", handleOrder)
	app.Get("/restricted/:uuid", handleRestricted)
	app.Use("/papi", handlePapi)
}
