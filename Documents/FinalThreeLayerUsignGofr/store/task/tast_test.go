package store

import (
	"FinalThreeLayerUsingGOfr/models"
	_ "database/sql"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"golang.org/x/net/context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAllTasks(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t) // a Fake Container for testing

	s := New()

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	expectedTasks := []models.Task{
		{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
		{ID: 2, Title: "Task 2", UserID: 1, Completed: true},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id", "completed"}).
		AddRow(expectedTasks[0].ID, expectedTasks[0].Title, expectedTasks[0].UserID, expectedTasks[0].Completed).
		AddRow(expectedTasks[1].ID, expectedTasks[1].Title, expectedTasks[1].UserID, expectedTasks[1].Completed)

	mock.SQL.ExpectQuery("SELECT id, title, user_id, completed FROM tasks").WillReturnRows(rows)

	tasks, err := s.GetAllTasks(ctx)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(tasks) != len(expectedTasks) {
		t.Errorf("expected %d tasks, got %d", len(expectedTasks), len(tasks))
	}

	for i, task := range tasks {
		if task != expectedTasks[i] {
			t.Errorf("task %d mismatch: got %+v, expected %+v", i, task, expectedTasks[i])
		}
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTaskByID(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	s := New()

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	testID := 1
	expectedTask := models.Task{
		ID:     testID,
		Title:  "Test Task",
		UserID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "user_id"}).
		AddRow(expectedTask.ID, expectedTask.Title, expectedTask.UserID)

	mock.SQL.ExpectQuery("SELECT id, title, user_id FROM tasks WHERE id = ?").
		WithArgs(testID).
		WillReturnRows(rows)

	task, err := s.GetTaskByID(ctx, testID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if task != expectedTask {
		t.Errorf("task mismatch: got %+v, expected %+v", task, expectedTask)
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	s := New()

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	testTask := models.Task{
		Title:  "New Task",
		UserID: 1,
	}

	mock.SQL.ExpectExec("INSERT INTO tasks(title, user_id) VALUES(?, ?)").
		WithArgs(testTask.Title, testTask.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := s.CreateTask(ctx, testTask)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	s := New()

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	testID := 1

	mock.SQL.ExpectExec(`UPDATE tasks SET completed = ? WHERE id = ?`).
		WithArgs(true, testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.UpdateTask(ctx, testID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	s := New()

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	testID := 1

	mock.SQL.ExpectExec("DELETE FROM tasks WHERE id = ?").
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.DeleteTask(ctx, testID)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

}
