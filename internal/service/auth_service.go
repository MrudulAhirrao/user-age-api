package service

import (
	"context"
	"time"
	"user-age-api/db/sqlc"
	"user-age-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"user-age-api/internal/websocket"
)


type AuthService struct{ 
	db *pgxpool.Pool
	queries *db.Queries
	jwtSecret string
	Hub *websocket.Hub
}

func NewAuthService(pool *pgxpool.Pool,secret string, hub *websocket.Hub) *AuthService{
	return &AuthService{
		db:	pool,        
        queries:   db.New(pool), // We can generate queries directly from the pool
        jwtSecret: secret,
		Hub: hub,
	}
}


func(s *AuthService) Login(ctx context.Context, req models.LoginRequest) (string, error){
	userDB, err := s.queries.GetUserByEmail(ctx,req.Email)
	if err != nil{
		return "",err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDB.PasswordHash),[]byte(req.Password))
	if err != nil{
		return "",err
	}

	claims:= jwt.MapClaims{
		"user_id": userDB.ID,
		"role": userDB.Role,
		"exp":	time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil{
		return "",err
	}
	return tokenString,nil
}


func (s *AuthService) Signup(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse,error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password),12)
	birthDate, err := time.Parse("2006-01-02", req.Dob)
	if err != nil{
		return nil, err
	}

	userDB, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Name: 	req.Name,
		Email:	req.Email,
		PasswordHash:	string(hashedPassword),
		Dob:          pgtype.Date{Time: birthDate, Valid: true},
		Role:	"user",
	})
	if err != nil{
		return nil, err
	}

	return &models.UserResponse{
		ID:	userDB.ID,
		Name: userDB.Name,
		Email:	userDB.Email,
		Dob:	userDB.Dob.Time.Format("2006-01-02"),
		Role:	userDB.Role,
		CreatedAt:	userDB.CreatedAt.Time,
	},nil
}

func (s *AuthService) GetMe(ctx context.Context, userID int32) (*models.UserResponse, error){
	userDB, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:	userDB.ID,
		Name:	userDB.Name,
		Email:	userDB.Email,
		Dob:	userDB.Dob.Time.Format("2006-01-02"),
		Age:  models.CalculateAge(userDB.Dob.Time),
		Role: userDB.Role,
		CreatedAt:	userDB.CreatedAt.Time,
	}, nil
}