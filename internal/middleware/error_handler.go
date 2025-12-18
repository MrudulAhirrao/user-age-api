package middleware
import(
	"errors"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx,err error) error{
	code := fiber.StatusInternalServerError
	message:= "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e){
		code = e.Code
		message = e.Message
	}
	return c.Status(code).JSON(fiber.Map{"error":true, "message": message,})
}