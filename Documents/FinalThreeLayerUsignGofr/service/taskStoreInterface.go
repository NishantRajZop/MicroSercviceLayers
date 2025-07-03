package service

import (
	"FinalThreeLayerUsingGOfr/models"
	"gofr.dev/pkg/gofr"
)

type TaskStore interface {
	GetAllTasks(ctx *gofr.Context) ([]models.Task, error)
	GetTaskByID(ctx *gofr.Context, id int) (models.Task, error)
	CreateTask(ctx *gofr.Context, task models.Task) (models.Task, error)
	UpdateTask(ctx *gofr.Context, id int) error
	DeleteTask(ctx *gofr.Context, id int) error
}
