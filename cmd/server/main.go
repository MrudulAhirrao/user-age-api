package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-age-api/config"
	"user-age-api/db/sqlc"
	"user-age-api/internal/handler"
	"user-age-api/internal/middleware"
	myredis "user-age-api/internal/redis"
	"user-age-api/internal/service"

	"user-age-api/internal/websocket" // ✅ Local Package (Your Hub)

	// ✅ External Package (Aliased to avoid conflict)
	fiberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

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

	hub := websocket.NewHub()
	go hub.Run()
	userService := service.NewUserService(queries)
	authService := service.NewAuthService(dbPool,cfg.JWTSecret,hub)
	userHandler := handler.NewUserHandler(userService, zapLogger)
	authHandler := handler.NewAuthHandler(authService,zapLogger)


	redisClient,err := myredis.NewRedisClient(cfg.RedisURL,zapLogger)

	if err != nil{
		zapLogger.Fatal("Could Not Connect to Redis",zap.Error(err))
	}
	defer redisClient.Close()


	
	healthHandler:= handler.NewHealthHandler(dbPool,redisClient)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler(zapLogger),
	})


	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.ZapLogger(zapLogger))
	app.Use(cors.New())

	app.Get("/healthz",healthHandler.Check)
	authapi := app.Group("/auth")
	authapi.Post("/signup", authHandler.Signup)
	authapi.Post("/login",middleware.RateLimitMiddleware(redisClient, 5, 60*time.Second),authHandler.Login)
	api := app.Group("/users",middleware.RateLimitMiddleware(redisClient, 100, 60*time.Second))
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	app.Get("/ws", func(c *fiber.Ctx) error {
		// ✅ FIX: Use 'fiberws' alias here
		if fiberws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}, handler.ServeWS(hub))

	api.Put("/:id/profile",authHandler.UpdateProfile)
	api.Get("/me",authHandler.GetMe)
	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.ListUsers)
	api.Get("/:id", userHandler.GetUser)
	api.Put("/:id", userHandler.UpdateUser)
	// api.Delete("/:id", userHandler.DeleteUser)

	//admin
	api.Delete("/:id", middleware.RoleMiddleware("admin"),userHandler.DeleteUser)
	
	c := make(chan os.Signal,1)
	signal.Notify(c, os.Interrupt,syscall.SIGTERM)

	go func(){
		if err := app.Listen(":3000"); err != nil{
			zapLogger.Info("Server was shutting Down")
		}
	}()

	zapLogger.Info("Server Running")
	<-c
	zapLogger.Info("Shutting down the server")

	ctx, cancel:= context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
        zapLogger.Error("Server forced to shutdown", zap.Error(err))
    }

    // 6. Close Infrastructure
    zapLogger.Info("Closing Database and Redis...")
    dbPool.Close()
    redisClient.Close()

    zapLogger.Info("👋 Server exited successfully")
	os.Exit(0)
	log.Fatal(app.Listen(":3000"))



}