package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"user-age-api/db/sqlc"
	emailClient "user-age-api/internal/client/email"
	"user-age-api/internal/models"
	"user-age-api/internal/websocket"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)


type AuthService struct{ 
	db *pgxpool.Pool
	queries *db.Queries
	jwtSecret string
	Hub *websocket.Hub
	emailClient *emailClient.EmailClinet
	redisClient *redis.Client
}

func NewAuthService(pool *pgxpool.Pool,secret string, hub *websocket.Hub, eClient *emailClient.EmailClinet, rClient *redis.Client) *AuthService{
	return &AuthService{
		db:	pool,        
        queries:   db.New(pool), // We can generate queries directly from the pool
        jwtSecret: secret,
		Hub: hub,
		emailClient: eClient,
		redisClient: rClient,
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}


func(s *AuthService) Login(ctx context.Context, req models.LoginRequest) (string, error){
	userDB, err := s.queries.GetUserByEmail(ctx,req.Email)
	if err != nil{
		return "",err
	}

	if !userDB.IsActive {
		return "", errors.New("account not activated. please check your email")
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

	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	userDB, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Name: 	req.Name,
		Email:	req.Email,
		PasswordHash:	string(hashedPassword),
		Dob:          pgtype.Date{Time: birthDate, Valid: true},
		Role:	"user",
		ActivationToken: pgtype.Text{String: token, Valid: true},
	})
	if err != nil{
		return nil, err
	}

	go func() {
		_ = s.emailClient.SendActivationEmail(userDB.Email, token)
	}()

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

func (s *AuthService) ActivateAccount(ctx context.Context, token string) error {
	// Try to find user with this token and set is_active = true
	// We assume pgtype.Text logic handles the nullable string
	_, err := s.queries.ActivateUser(ctx, pgtype.Text{String: token, Valid: true})
	if err != nil {
		return errors.New("invalid or expired activation token")
	}
	return nil
}

func (s *AuthService) ResendActivationEmail(ctx context.Context, email string) error {
	limitKey := "resend_limit:" + email      
	cooldownKey := "resend_cooldown:" + email 
	if s.redisClient.Exists(ctx, cooldownKey).Val() > 0 {
		ttl := s.redisClient.TTL(ctx, cooldownKey).Val()
		return errors.New("please wait " + ttl.Round(time.Second).String() + " before retrying")
	}
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}
	if user.IsActive {
		return errors.New("account already active")
	}

	count := s.redisClient.Incr(ctx, limitKey).Val()
	if count == 1 {
		s.redisClient.Expire(ctx, limitKey, 24*time.Hour) 
	}

	var waitTime time.Duration
	switch count {
	case 1:
		waitTime = 1 * time.Minute
	case 2:
		waitTime = 5 * time.Minute
	case 3:
		waitTime = 20 * time.Minute
	default:
		waitTime = 24 * time.Hour 
	}

	go func() {
		_ = s.emailClient.SendActivationEmail(user.Email, user.ActivationToken.String)
	}()

	s.redisClient.Set(ctx, cooldownKey, "blocked", waitTime)

	return nil
}