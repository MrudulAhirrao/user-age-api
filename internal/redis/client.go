package redis

import (
	"context"
	// "errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct{
	Client *redis.Client
}

func NewRedisClient(url string, logger *zap.Logger) (*RedisClient,error){
	opt,err:=redis.ParseURL(url)
	if err!=nil{
		return nil, fmt.Errorf("Invalid Redis Url: %w",err)
	}
	client:= redis.NewClient(opt)

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	if err:= client.Ping(ctx).Err(); err!= nil{
		return nil,fmt.Errorf("Failed to connect%w",err)
	}
	logger.Info("Conncetion Succefuully to Redis")
	return &RedisClient{Client: client},nil
}

func(r *RedisClient) Close() error{
	return r.Client.Close()
}