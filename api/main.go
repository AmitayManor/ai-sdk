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
	"time"
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

	// init server engine
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	userHandler := handlers.NewUserHandler(config.GetSupabaseClient())
	authHandler := handlers.NewAuthHandler(config.GetSupabaseClient())

	//SignUp route
	app.Post("api/auth/signup",
		middleware.ValidateSignUp(),
		authHandler.SignUp,
	)

	//Public routes
	app.Post("/api/users", userHandler.CreateUser)
	app.Post("/api/auth/signin", middleware.RateLimiter(5, time.Minute),
		userHandler.SignIn,
	)

	//Protected routes
	api := app.Group("/api", middleware.Protected(),
		middleware.RateLimiter(100, time.Minute),
	)

	api.Get("/users", userHandler.ListUsers)
	api.Get("/users/:id", userHandler.GetUser)
	api.Put("/users/:id", userHandler.UpdateUser)
	api.Delete("/users/:id", userHandler.DeleteUser)
	api.Put("/users/:id/login-attempts", userHandler.UpdateLoginAttempts)
	api.Put("/users/:id/reset-attempts", userHandler.ResetLoginAttempts)

	//Admin routes
	admin := api.Group("/admin", middleware.AdminOnly(),
		middleware.RateLimiter(50, time.Minute),
	)
	admin.Put("/users/:id", userHandler.AdminUpdateUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
