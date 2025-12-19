package handler

import (
	"context"
	"time"
	"user-age-api/internal/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type healthHandler struct{
	DB *pgxpool.Pool
	Redis *redis.RedisClient
}

func NewHealthHandler(db *pgxpool.Pool, redis *redis.RedisClient) *healthHandler{
	return &healthHandler{
		DB:db, Redis: redis,
	}
}

func (h *healthHandler) Check(c* fiber.Ctx) error{
	ctx, cancel:= context.WithTimeout(context.Background(),2*time.Second)
	defer cancel()
	if err := h.DB.Ping(ctx); err != nil{
		return c.Status(503).JSON(fiber.Map{"status":"error","error":"DB Unavailable",})
	}
	if err := h.Redis.Client.Ping(ctx).Err(); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "error",
			"error":  "Redis unavailable",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status":"ok",
		"components": fiber.Map{
			"database":"Database Connected",
			"redis":"Redis Connected",
		},
	})
}