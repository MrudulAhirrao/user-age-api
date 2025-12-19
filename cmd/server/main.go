package main

import (
	"context"
	"log"
	"time"
	myredis "user-age-api/internal/redis"
	"user-age-api/config"
	"user-age-api/db/sqlc"
	"user-age-api/internal/handler"
	"user-age-api/internal/service"
	"user-age-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	// Initialize Zap Logger
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	dbPool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		zapLogger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer dbPool.Close()
	queries := db.New(dbPool)
	userService := service.NewUserService(queries)
	authService := service.NewAuthService(dbPool,cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService, zapLogger)
	authHandler := handler.NewAuthHandler(authService,zapLogger)


	redisClient,err := myredis.NewRedisClient(cfg.RedisURL,zapLogger)

	if err != nil{
		zapLogger.Fatal("Could Not Connect to Redis",zap.Error(err))
	}
	defer redisClient.Close()


	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})



	app.Use(middleware.RequestIDMiddleware())
	app.Use(logger.New(logger.Config{
		Format:"[${time}] ${locals:request_id} ${status} - ${method} ${path}\n",
	}))

	
	authapi := app.Group("/auth")
	authapi.Post("/signup", authHandler.Signup)
	authapi.Post("/login",middleware.RateLimitMiddleware(redisClient, 5, 60*time.Second),authHandler.Login)
	api := app.Group("/users")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	api.Put("/:id/profile",authHandler.UpdateProfile)
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