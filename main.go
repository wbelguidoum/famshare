package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	// Static files
	app.Static("/", "./static")

	// Public routes
	app.Post("/api/login", loginHandler)
	app.Post("/api/logout", logoutHandler)

	// Protected routes
	protected := app.Group("/api", authMiddleware)
	protected.Get("/session/check", checkSessionHandler)
	protected.Get("/photos", listPhotosHandler)
	protected.Post("/photos", uploadPhotoHandler)
	protected.Get("/photos/:id", getPhotoHandler)
	protected.Put("/photos/:id", updatePhotoHandler)
	protected.Delete("/photos/:id", deletePhotoHandler)

	log.Println("Starting FamShare v0.1 on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
