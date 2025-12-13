package models

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	// 1. Setup specific date for consistency in tests
	// Assume "Today" is fixed for the test, or calculate relative to Now
	
	tests := []struct {
		name     string
		dob      string // YYYY-MM-DD
		expected int
	}{
		{
			name:     "Born 2000 (Birthday passed)",
			dob:      time.Now().AddDate(-25, -1, 0).Format("2006-01-02"), // 25 years and 1 month ago
			expected: 25,
		},
		{
			name:     "Born 2000 (Birthday tomorrow)",
			dob:      time.Now().AddDate(-25, 0, 1).Format("2006-01-02"), // 25 years minus 1 day ago
			expected: 24, // Hasn't turned 25 yet
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedDob, _ := time.Parse("2006-01-02", tt.dob)
			age := CalculateAge(parsedDob)
			if age != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, age)
			}
		})
	}
}