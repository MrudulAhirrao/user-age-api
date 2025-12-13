package main

import (
	"context"
	"log"
	

	"user-age-api/db/sqlc"
	"user-age-api/internal/handler"
	"user-age-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Zap Logger
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	// 2. Connect to Database
	// Note: In production, use environment variables!
	// Docker connection string: postgres://user:secret@localhost:5432/userdb?sslmode=disable
	connStr := "postgres://user:secret@localhost:5432/userdb?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		zapLogger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer conn.Close(context.Background())

	// 3. Dependency Injection
	queries := db.New(conn)
	userService := service.NewUserService(queries)
	userHandler := handler.NewUserHandler(userService, zapLogger)

	// 4. Setup Fiber
	app := fiber.New()

	// Add Middleware (Logger + Request ID is automatic in Fiber logs usually, but let's add basic logging)
	app.Use(logger.New())

	// 5. Define Routes
	api := app.Group("/users")
	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.ListUsers)
	api.Get("/:id", userHandler.GetUser)
	api.Put("/:id", userHandler.UpdateUser)
	api.Delete("/:id", userHandler.DeleteUser)

	// 6. Start Server
	log.Fatal(app.Listen(":3000"))
}