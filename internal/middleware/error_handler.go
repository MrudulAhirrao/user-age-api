package middleware

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GlobalErrorHandler formats all errors into a standard JSON response
func GlobalErrorHandler(log *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// 1. Default Defaults (500 Internal Server Error)
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		// 2. Check if it's a specific Fiber Error (e.g., 404 Not Found, 429 Too Many Requests)
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			message = e.Message
		}

		requestId := c.Locals("requestid") 
		
		if code >= 500 {
			log.Error("Server Error", 
				zap.Error(err), 
				zap.String("path", c.Path()),
				zap.Any("request_id", requestId),
			)
		} else {
			log.Info("Client Error", 
				zap.Int("code", code), 
				zap.String("error", message),
				zap.Any("request_id", requestId),
			)
		}

		// 4. Return the Clean JSON
		return c.Status(code).JSON(fiber.Map{
			"status":     "error",
			"message":    message,
			"requestId":  requestId,
		})
	}
}