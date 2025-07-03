package service

import (
	"FinalThreeLayerUsingGOfr/models"
	"errors"
	"gofr.dev/pkg/gofr"
)

type taskService struct {
	store TaskStore
}

func NewTaskService(store TaskStore) *taskService {
	return &taskService{store: store}
}

func (tService *taskService) GetAllTasks(ctx *gofr.Context) ([]models.Task, error) {
	//You can Perform All Your Validations here
	return tService.store.GetAllTasks(ctx)
}

func (tService *taskService) GetTaskByID(ctx *gofr.Context, id int) (models.Task, error) {
	//You can Perform All Your Validations here
	return tService.store.GetTaskByID(ctx, id)
}

func (tService *taskService) CreateTask(ctx *gofr.Context, task models.Task) (models.Task, error) {
	if task.Title == "" {
		return models.Task{}, errors.New("task title cannot be empty")
	}
	if task.UserID == 0 {
		return models.Task{}, errors.New("invalid user ID")
	}

	return tService.store.CreateTask(ctx, task)
}

func (tService *taskService) UpdateTask(ctx *gofr.Context, id int) error {
	if id == 0 {
		return errors.New("invalid task ID")
	}
	return tService.store.UpdateTask(ctx, id)
}

func (tService *taskService) DeleteTask(ctx *gofr.Context, id int) error {
	//You can Perform All Your Validations here
	if id == 0 {
		return errors.New("invalid task ID")
	}
	return tService.store.DeleteTask(ctx, id)
}
