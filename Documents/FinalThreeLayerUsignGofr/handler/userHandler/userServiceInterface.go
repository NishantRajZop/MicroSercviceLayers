package userHandler

import (
	"FinalThreeLayerUsingGOfr/models"
	"gofr.dev/pkg/gofr"
)

type userService interface {
	CreateUser(*gofr.Context, models.User) (models.User, error)
	GetUserByID(ctx *gofr.Context, id int) (models.User, error)
}
