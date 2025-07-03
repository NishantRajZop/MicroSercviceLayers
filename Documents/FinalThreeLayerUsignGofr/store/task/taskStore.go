package store

import (
	"FinalThreeLayerUsingGOfr/models"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
)

type taskStore struct {
	db *sql.DB
}

// this new Function is like  a constructor  , if You will call this new with a db then it will return you a taskStore instance
// which will have all the CRUD functionalites🥳

// usually this is called by a service layer , so that service layer can have a kind of stick that can operate CRUD

func New() *taskStore {
	return &taskStore{}
}

func (s *taskStore) GetAllTasks(ctx *gofr.Context) ([]models.Task, error) {
	rows, err := ctx.SQL.Query("SELECT id, title, user_id, completed FROM tasks")
	defer rows.Close()

	if err != nil {
		return []models.Task{}, errors.New("database error")
	}

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		rows.Scan(&t.ID, &t.Title, &t.UserID, &t.Completed)
		tasks = append(tasks, t)
	}
	//fmt.Println(tasks)
	return tasks, nil
}

func (s *taskStore) GetTaskByID(ctx *gofr.Context, id int) (models.Task, error) {
	var task models.Task
	err := ctx.SQL.QueryRow("SELECT id, title, user_id FROM tasks WHERE id = ?", id).
		Scan(&task.ID, &task.Title, &task.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Task{}, fmt.Errorf("task with id %d not found", id) // or custom NotFound error
		}
		return models.Task{}, err // generic DB error
	}

	return task, nil
}

func (s *taskStore) CreateTask(ctx *gofr.Context, task models.Task) (models.Task, error) {
	result, err := ctx.SQL.Exec("INSERT INTO tasks(title, user_id) VALUES(?, ?)", task.Title, task.UserID)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to insert task: %v", err)
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to retrieve inserted ID: %v", err)
	}

	// Build the task with the new ID
	createdTask := models.Task{
		ID:     int(id),
		Title:  task.Title,
		UserID: task.UserID,
	}

	return createdTask, nil
}

func (s *taskStore) UpdateTask(ctx *gofr.Context, id int) error {
	result, err := ctx.SQL.Exec("UPDATE tasks SET completed = ? WHERE id = ?", true, id)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return gofrHttp.ErrorEntityNotFound{Name: fmt.Sprintf("Task ID %d", id)}
	}

	return nil
}

func (s *taskStore) DeleteTask(ctx *gofr.Context, id int) error {
	result, err := ctx.SQL.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deletion result: %v", err)
	}

	if rowsAffected == 0 {
		// Return a structured not-found error
		return &gofrHttp.ErrorEntityNotFound{
			Name:  "task",
			Value: fmt.Sprintf("%d", id),
		}
		// Alternatively, simple error:
		// return fmt.Errorf("task with id %d not found", id)
	}
	return nil
}
