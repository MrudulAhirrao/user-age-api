package middleware

import (
	"net/http"
	"time"
	"fmt"
	// "github.com/redis/go-redis/v9"
	"github.com/gofiber/fiber/v2"
	myredis "user-age-api/internal/redis"
)
func RateLimitMiddleware(r *myredis.RedisClient,limit int, window time.Duration) fiber.Handler{
	return func(c *fiber.Ctx) error{
		ip:= c.IP()
		path:= c.Path()

		key:= fmt.Sprintf("rate_limit:%s%s",ip,path)

		pipe:= r.Client.Pipeline()

		incr:= pipe.Incr(c.Context(),key)
		pipe.Expire(c.Context(),key,window)

		_,err:= pipe.Exec(c.Context())
		if err!=nil{
			fmt.Println("⚠️ DEBUG: Redis failed:", err)
			return c.Next()
		}
		count:= incr.Val()
		c.Set("X-Rate-Limit",fmt.Sprintf("%d",limit))
		c.Set("X-Rate-Limit-Remaining",fmt.Sprintf("%d",int64(limit)-count))
		fmt.Println("🔍 DEBUG: Count for", ip, "is", count)
		if count > int64(limit){
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{"error":"Too Many Request. Please Try Again Later","retry_after":window.Seconds() })
		}
		return c.Next()
	}
}