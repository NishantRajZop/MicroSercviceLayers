package store

import (
	"FinalThreeLayerUsingGOfr/models"
	"github.com/DATA-DOG/go-sqlmock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"golang.org/x/net/context"
	"testing"
)

func TestCreateUser(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	s := New()

	testUser := models.User{
		ID:   1,
		Name: "test_user1 1",
	}

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	// Creates a new instance of your repository/struct (presumably named str)
	// Passes the mock database connection to it
	mock.SQL.ExpectExec("INSERT INTO users(name) VALUES(?)").WithArgs(testUser.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := s.CreateUser(ctx, testUser)

	if err != nil {
		t.Errorf("expected success message, got: %s", err)
	}
}

func TestGetUserByID(t *testing.T) {

	mockContainer, mock := container.NewMockContainer(t)

	s := New()

	testID := 1
	expectedUser := models.User{
		ID:   testID,
		Name: "test_user1",
	}

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedUser.ID, expectedUser.Name)

	mock.SQL.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(testID).
		WillReturnRows(rows)

	user, err := s.GetUserByID(ctx, testID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if user.ID != expectedUser.ID {
		t.Errorf("got user ID %d, expected %d", user.ID, expectedUser.ID)
	}
	if user.Name != expectedUser.Name {
		t.Errorf("got user name '%s', expected '%s'", user.Name, expectedUser.Name)
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
