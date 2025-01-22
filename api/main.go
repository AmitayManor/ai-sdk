package main

import (
	"api/config"
	"api/handlers"
	"api/middleware"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

func main() {

	// init supabase client
	if err := config.InitSupabase(); err != nil {
		log.Fatal("Failed to initialized Supabase: %v", err)

	}

	//init project db
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	//// Example query to test connection
	//var version string
	//if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
	//	log.Fatalf("Query failed: %v", err)
	//}
	//
	//log.Println("Connected to:", version)
	//
	// ITS WORKING!!!!!!!!!

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

	//app.Get("/api/test/users", userHandler.TestConnection)

	//Public routes
	app.Post("/api/users", userHandler.CreateUser)
	app.Post("/api/auth/signin", userHandler.SignIn)

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
