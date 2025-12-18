package main

import (
	"context"
	"log"
	
	"user-age-api/config"
	"user-age-api/db/sqlc"
	"user-age-api/internal/handler"
	"user-age-api/internal/service"
	"user-age-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	// Initialize Zap Logger
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	conn, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		zapLogger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer conn.Close(context.Background())
	queries := db.New(conn)
	userService := service.NewUserService(queries)
	authService := service.NewAuthService(queries,cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService, zapLogger)
	authHandler := handler.NewAuthHandler(authService,zapLogger)


	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(middleware.RequestIDMiddleware())
	app.Use(logger.New(logger.Config{
		Format:"[${time}] ${locals:request_id} ${status} - ${method} ${path}\n",
	}))

	authapi := app.Group("/auth")
	authapi.Post("/signup", authHandler.Signup)
	authapi.Post("/login",authHandler.Login)
	api := app.Group("/users")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	api.Get("/me",authHandler.GetMe)
	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.ListUsers)
	api.Get("/:id", userHandler.GetUser)
	api.Put("/:id", userHandler.UpdateUser)
	// api.Delete("/:id", userHandler.DeleteUser)

	//admin
	api.Delete("/:id", middleware.RoleMiddleware("admin"),userHandler.DeleteUser)
	
	log.Fatal(app.Listen(":3000"))
}