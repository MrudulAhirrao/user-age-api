package models

import "time"

// CreateUserRequest is the JSON body for POST/PUT
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"` // Enforce format YYYY-MM-DD
}

// UserResponse is what we send back to the client
type UserResponse struct {
	ID   int32     `json:"id"`
	Name string    `json:"name"`
	Dob  string    `json:"dob"`
	Age  int       `json:"age"` // Age is calculated dynamically
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