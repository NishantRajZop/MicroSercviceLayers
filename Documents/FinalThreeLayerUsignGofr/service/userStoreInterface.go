package service

import (
	"FinalThreeLayerUsingGOfr/models"
	"gofr.dev/pkg/gofr"
)

type UserStore interface {
	CreateUser(*gofr.Context, models.User) (models.User, error)
	GetUserByID(*gofr.Context, int) (models.User, error)
}
