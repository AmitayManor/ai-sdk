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
		log.Fatalf("Failed to initialized Supabase: %v", err)

	}

	//init project db
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	//init postgres client
	if err := config.InitPostgres(); err != nil {
		log.Fatalf("Failed to initialized Postgres Client: %v", err)
	}

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

	authHandler := handlers.NewAuthHandler(config.GetSupabaseClient())
	storageHandler := handlers.NewStorageHandler(config.GetStorageClient())
	requestHandler := handlers.NewRequestHandler()

	//SignUp route
	app.Post("api/auth/signup",
		middleware.ValidateSignUp(),
		authHandler.SignUp,
	)

	app.Post("/auth/signin", authHandler.SignIn)

	// Protected routes
	api := app.Group("/api", middleware.AuthRequired())

	// Model requests endpoints
	api.Post("/requests", requestHandler.CreateModelRequest)
	api.Get("/requests", requestHandler.ListRequests)

	// Storage endpoints
	api.Get("/images/:id", storageHandler.GetGeneratedImage)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
