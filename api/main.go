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

	middleware.InitBlacklist()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	userHandler := handlers.NewUserHandler(config.GetSupabaseClient())
	authHandler := handlers.NewAuthHandler(config.GetSupabaseClient())
	apiKeyHandler := handlers.NewAPIKeyHandler()
	modelHandler := handlers.NewModelHandler()
	requestHandler := handlers.NewRequestHandler()

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
	api.Post("/requests", middleware.RateLimiter(50, time.Minute), requestHandler.CreateRequest)
	api.Get("/requests", middleware.RateLimiter(100, time.Minute), requestHandler.ListRequests)
	api.Get("/requests/:id", middleware.RateLimiter(100, time.Minute), requestHandler.GetRequest)

	//Admin routes
	admin := api.Group("/admin", middleware.AdminOnly(),
		middleware.RateLimiter(50, time.Minute),
	)
	admin.Put("/users/:id", userHandler.AdminUpdateUser)
	admin.Post("/models", middleware.RateLimiter(20, time.Minute), modelHandler.CreateModel)
	admin.Get("/models", middleware.RateLimiter(100, time.Minute), modelHandler.ListModels)
	admin.Get("/models/:id", middleware.RateLimiter(100, time.Minute), modelHandler.GetModel)
	admin.Put("/models/:id", middleware.RateLimiter(20, time.Minute), modelHandler.UpdateModel)
	admin.Delete("/models/:id", middleware.RateLimiter(20, time.Minute), modelHandler.DeleteModel)

	keys := api.Group("/keys")
	keys.Post("/", apiKeyHandler.CreateKey)
	keys.Get("/", apiKeyHandler.ListKeys)
	keys.Delete("/:id", apiKeyHandler.DeactivateKey)
	keys.Put("/:id", apiKeyHandler.UpdateKey)

	keyProtected := app.Group("/api/v1", middleware.ValidateAPIKey())
	keyProtected.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
