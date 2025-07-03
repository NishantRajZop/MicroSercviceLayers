package store

import (
	"FinalThreeLayerUsingGOfr/models"
	"database/sql"
	"fmt"
	"gofr.dev/pkg/gofr"
)

type userStore struct {
	db *sql.DB
}

func New() *userStore {
	return &userStore{}
}

func (s *userStore) CreateUser(ctx *gofr.Context, user models.User) (models.User, error) {
	result, err := ctx.SQL.Exec("INSERT INTO users(name) VALUES(?)", user.Name)
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve inserted ID: %v", err)
	}

	// Return the created user with the new ID
	return models.User{
		ID:   int(id),
		Name: user.Name,
	}, nil
}

func (s *userStore) GetUserByID(ctx *gofr.Context, id int) (models.User, error) {
	var user models.User
	err := ctx.SQL.QueryRow("SELECT id, name FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name)
	if err != nil {
		return models.User{}, fmt.Errorf("database error: %v", err)
	}

	return user, nil
}
