package models

import "time"

// CreateUserRequest is the JSON body for POST/PUT
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate: "required,email"`
	Password string `json: "password" validate: "required, min=8"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}
type LoginRequest struct{
	Email string	`json: "email" validate: "required, email"`
	Password string `json: "password" validate: "required`
}
// UserResponse is what we send back to the client
type UserResponse struct {
	ID   int32     `json:"id"`
	Name string    `json:"name"`
	Email string   `json:"email"`
	Dob  string    `json:"dob"`
	Age  int       `json:"age"` 
	Role string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Helper to calculate age
func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	
	// If current date is before birthday this year, subtract 1
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}