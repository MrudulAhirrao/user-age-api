package middleware

import "github.com/gofiber/fiber/v2"

func RoleMiddleware(requiredRole string) fiber.Handler{
	return func(c *fiber.Ctx) error{
		userRole, ok:= c.Locals("role").(string)

		if !ok || userRole == ""{
			return c.Status(403).JSON(fiber.Map{"error":"Frobidden Error- Role not found"})
		}
		if userRole != requiredRole{
			return c.Status(403).JSON(fiber.Map{"error":"Forbidden error- You don't have permission"})
		}
		return c.Next()
	}
}