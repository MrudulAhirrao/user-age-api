package service

import (
	"context"
	"time"
	db "user-age-api/db/sqlc"
	"user-age-api/internal/models"

	"github.com/jackc/pgx/v5/pgtype" 
)

type UserService struct {
	queries *db.Queries
}

func NewUserService(queries *db.Queries) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	parsedDob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Name: req.Name,
		Dob:  pgtype.Date{Time: parsedDob, Valid: true}, 
	})
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  models.CalculateAge(user.Dob.Time),
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, id int32) (models.UserResponse, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  models.CalculateAge(user.Dob.Time),
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	var response []models.UserResponse
	for _, u := range users {
		response = append(response, models.UserResponse{
			ID:   u.ID,
			Name: u.Name,
			Dob:  u.Dob.Time.Format("2006-01-02"),
			Age:  models.CalculateAge(u.Dob.Time),
		})
	}
	// Return empty slice instead of nil if no users found (better JSON)
	if response == nil {
		response = []models.UserResponse{}
	}
	return response, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.CreateUserRequest) (models.UserResponse, error) {
	parsedDob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	
	user, err := s.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob:  pgtype.Date{Time: parsedDob, Valid: true},
	})
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  models.CalculateAge(user.Dob.Time),
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.queries.DeleteUser(ctx, id)
}