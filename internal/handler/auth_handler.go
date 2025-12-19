package handler

import(
	"user-age-api/internal/models"
	"user-age-api/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct{
	service	*service.AuthService
	logger	*zap.Logger
	validator	*validator.Validate
}

func NewAuthHandler(s *service.AuthService, l *zap.Logger) *AuthHandler{
	return &AuthHandler{
		service: s,
		logger: l,
		validator:	validator.New(),
	}
}


func (h* AuthHandler) Login(c * fiber.Ctx) error{
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil{
		return  c.Status(400).JSON(fiber.Map{"error":"Invalid Json"})
	}
	token, err :=h.service.Login(c.Context(),req)
	if err!=nil{
		
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Email or Password"})
	}
	return c.Status(200).JSON(fiber.Map{
		"token":token,
		"type": "Bearer",
	})
}


func (h *AuthHandler) Signup(c *fiber.Ctx) error{
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	if err := h.validator.Struct(req); err!= nil{
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	//Service Called here 
	user, err := h.service.Signup(c.Context(), req)
	if err != nil{
		h.logger.Error("Signup Failed", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": "Could not Register the User"})
	}
	return c.Status(201).JSON(user)
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error{
	userID, ok := c.Locals("user_id").(int32)
	if !ok{
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorizeds: No User Found"})
	}
	user, err := h.service.GetMe(c.Context(), userID)
	if err != nil{
		h.logger.Error("Profile Fetching Failed", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error":"Could not Fetch the Profile"})
	}
	return c.Status(200).JSON(user)
}

func (h*AuthHandler) UpdateProfile(c *fiber.Ctx) error{
	paramID, err:= c.ParamsInt("id")
	if err != nil{
		return c.Status(400).JSON(fiber.Map{"error":"Invalid User ID"})
	}
	loggedUserID := c.Locals("user_id").(int32)

	if int32(paramID) != loggedUserID{
		return c.Status(403).JSON(fiber.Map{"error":"You Can Update Only Your Profile"})
	}

	var req service.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Request Body"})
	}

	req.UserID = int32(paramID)

	err = h.service.UpdateProfile(c.Context(),req)
	if err != nil{
		h.logger.Error("Failed to Update the Profile", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error":"Failed to update profile"})
	}
	return c.Status(200).JSON(fiber.Map{"message":"Profile Updated Successfully."})
}