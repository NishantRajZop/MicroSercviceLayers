package service

import (
	"FinalThreeLayerUsingGOfr/models"
	"errors"
	"gofr.dev/pkg/gofr"
)

type userService struct {
	store UserStore
}

func NewUserService(store UserStore) *userService {
	return &userService{store: store}
}

func (uService *userService) CreateUser(ctx *gofr.Context, user models.User) (models.User, error) {
	// you can implement all your Validations here
	if user.Name == "" {
		return models.User{}, errors.New("user name cannot be empty")
	}
	if user.ID == 0 {
		return models.User{}, errors.New("invalid user ID")
	}
	return uService.store.CreateUser(ctx, user)
}

func (uService *userService) GetUserByID(ctx *gofr.Context, id int) (models.User, error) {
	if id == 0 {
		return models.User{}, errors.New("invalid user ID")
	}
	return uService.store.GetUserByID(ctx, id)
}
