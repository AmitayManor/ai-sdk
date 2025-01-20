package main

import (
	"api/config"
	"api/handlers"
	"api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {

	if err := config.InitSupabase(); err != nil {
		log.Fatal("Failed to initialized Supabase: %v", err)

	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())

	userHandler := handlers.NewUserHandler(config.GetSupabaseClient())

	//Public routes
	app.Post("/api/users", userHandler.CreateUser)

	//Protected routes
	api := app.Group("/api", middleware.Protected())
	api.Get("/users", userHandler.ListUsers)
	api.Get("/users/:id", userHandler.GetUser)
	api.Put("/users/:id", userHandler.UpdateUser)
	api.Delete("/users/:id", userHandler.DeleteUser)
	api.Put("/users/:id/login-attempts", userHandler.UpdateLoginAttempts)
	api.Put("/users/:id/reset-attempts", userHandler.ResetLoginAttempts)

	//Admin routes
	admin := api.Group("/admin", middleware.AdminOnly())
	admin.Put("/users/:id", userHandler.AdminUpdateUser)

	//Start server
	log.Fatal(app.Listen(":3000"))
}
