package models

type Task struct {
	ID        int    `json:"id" example:"1"`
	Title     string `json:"title" example:"Complete project"`
	UserID    int    `json:"user_id" example:"123"`
	Completed bool   `json:"completed" example:"false"`
}
