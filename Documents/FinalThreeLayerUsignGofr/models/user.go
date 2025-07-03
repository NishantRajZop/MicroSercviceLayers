package models

type User struct {
	ID   int    `json:"id" example:"123"`        // User ID example
	Name string `json:"name" example:"John Doe"` // User name example
}
