package middleware

import(
	"fmt"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) fiber.Handler{
	return func(c *fiber.Ctx)error{
		authHeader := c.Get("Authorization")
		if authHeader == ""{
			return c.Status(400).JSON(fiber.Map{"error":"Missing Authrization Header"})
		}
		parts := strings.Split(authHeader," ")
		if len(parts) != 2 || parts[0] != "Bearer"{
			return c.Status(401).JSON(fiber.Map{"error":"Invalid Header Format"})
		}
		tokenString:= parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{},error){
			if _, ok:= token.Method.(*jwt.SigningMethodHMAC); !ok{
				return nil, fmt.Errorf("unexpected siging method: %v", token.Header["alg"])
			}
			return []byte(secret),nil
		})
		if err != nil || !token.Valid{
			return c.Status(401).JSON(fiber.Map{"error":"Invalid / Expired Token"})
		}
		claims,ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid{
			userID := int32(claims["user_id"].(float64))
			role:= claims["role"].(string)

			c.Locals("user_id",userID)
			c.Locals("role",role)
		}
		return c.Next()
	}
}