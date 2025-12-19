package middleware

import (
	"time"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ZapLogger(log *zap.Logger) fiber.Handler{
	return func(c *fiber.Ctx) error{
		start:= time.Now()
		err:= c.Next()
		duration:= time.Since(start)

		statusCode:= c.Response().StatusCode()

		requestID:= c.Locals("requestid")

		if err != nil {
			var e *fiber.Error
			if errors.As(err, &e) {
				statusCode = e.Code
			} else {
				// If it's a generic Go error, default to 500
				statusCode = 500
			}
		}

		fields:=[]zap.Field{
			zap.Int("status",statusCode),
			zap.String("method",c.Method()),
			zap.String("path",c.Path()),
			zap.String("ip",c.IP()),
			zap.Duration("latency",duration),
			zap.Any("req_id",requestID),
		}
		if err != nil{
			fields = append(fields, zap.Error(err))
			log.Error("Request Failed",fields...)
		}else if statusCode >= 500 {
			log.Error("Server Error", fields...)
		}else if statusCode >= 400{
			if err != nil {
				fields = append(fields, zap.String("error_msg", err.Error()))
			}
			log.Warn("Client Error", fields...)
		}else{
			log.Info("Request Proceeds", fields...)
		}
		return err
	}
}