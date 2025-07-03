package service

import (
	"FinalThreeLayerUsingGOfr/models"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"golang.org/x/net/context"
	"testing"
)

func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockTaskStore(ctrl)

	expectedTasks := []models.Task{
		{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
		{ID: 2, Title: "Task 2", UserID: 1, Completed: true},
	}

	tests := []struct {
		name      string
		mockSetup func(*gofr.Context)
		wantTasks []models.Task
		wantErr   error
	}{
		{
			name: "Success",
			mockSetup: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					GetAllTasks(ctx).
					Return(expectedTasks, nil)
			},
			wantTasks: expectedTasks,
			wantErr:   nil,
		},
		{
			name: "Failure",
			mockSetup: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					GetAllTasks(ctx).
					Return(nil, errors.New("db failure"))
			},
			wantTasks: nil,
			wantErr:   errors.New("db failure"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := &gofr.Context{
				Context: context.Background(),
			}

			tt.mockSetup(ctx)

			s := NewTaskService(mockStore)
			gotTasks, err := s.GetAllTasks(ctx)

			assert.Equal(t, tt.wantTasks, gotTasks, "TEST[%d] - %s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d] - %s", i, tt.name)
		})
	}
}
func TestGetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockTaskStore(ctrl)

	tests := []struct {
		name             string
		taskID           int
		mockExpect       func(ctx *gofr.Context)
		expectedResponse models.Task
		expectedErr      error
	}{
		{
			name:   "Success - Get task by ID",
			taskID: 1,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					GetTaskByID(ctx, 1).
					Return(models.Task{ID: 1, Title: "Task 1", UserID: 1}, nil)
			},
			expectedResponse: models.Task{ID: 1, Title: "Task 1", UserID: 1},
			expectedErr:      nil,
		},
		{
			name:   "Error - Task not found",
			taskID: 999,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					GetTaskByID(ctx, 999).
					Return(models.Task{}, sql.ErrNoRows)
			},
			expectedResponse: models.Task{},
			expectedErr:      sql.ErrNoRows,
		},
		{
			name:   "Error - DB connection failed",
			taskID: 1,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					GetTaskByID(ctx, 1).
					Return(models.Task{}, errors.New("connection failed"))
			},
			expectedResponse: models.Task{},
			expectedErr:      errors.New("connection failed"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock context
			ctx := &gofr.Context{
				Context: context.Background(),
			}

			// Setup mock expectations
			tt.mockExpect(ctx)

			// Inject mockStore into service
			s := NewTaskService(mockStore)

			// Call service method
			gotTask, err := s.GetTaskByID(ctx, tt.taskID)

			// Assert result
			assert.Equal(t, tt.expectedResponse, gotTask, "TEST[%d] - %s", i, tt.name)
			assert.Equal(t, tt.expectedErr, err, "TEST[%d] - %s", i, tt.name)
		})
	}
}

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockTaskStore(ctrl)

	tests := []struct {
		name           string
		task           models.Task
		mockExpect     func(ctx *gofr.Context)
		expectedResult models.Task
		expectedError  error
	}{
		{
			name: "Success - Create task",
			task: models.Task{Title: "New Task", UserID: 1},
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					CreateTask(ctx, models.Task{Title: "New Task", UserID: 1}).
					Return(models.Task{ID: 1, Title: "New Task", UserID: 1, Completed: false}, nil)
			},
			expectedResult: models.Task{ID: 1, Title: "New Task", UserID: 1, Completed: false},
			expectedError:  nil,
		},
		{
			name: "Error - Empty title",
			task: models.Task{Title: "", UserID: 1},
			mockExpect: func(ctx *gofr.Context) {
				// No store call expected
			},
			expectedResult: models.Task{},
			expectedError:  errors.New("task title cannot be empty"),
		},
		{
			name: "Error - Invalid user ID",
			task: models.Task{Title: "New Task", UserID: 0},
			mockExpect: func(ctx *gofr.Context) {
				// No store call expected
			},
			expectedResult: models.Task{},
			expectedError:  errors.New("invalid user ID"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &gofr.Context{
				Context: context.Background(),
			}

			if tt.mockExpect != nil {
				tt.mockExpect(ctx)
			}

			s := NewTaskService(mockStore)

			result, err := s.CreateTask(ctx, tt.task)

			assert.Equal(t, tt.expectedResult, result, "TEST[%d] - %s", i, tt.name)
			assert.Equal(t, tt.expectedError, err, "TEST[%d] - %s", i, tt.name)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockTaskStore(ctrl)

	tests := []struct {
		name             string
		taskID           int
		mockExpect       func(ctx *gofr.Context)
		expectedResponse error
	}{
		{
			name:   "Success - Update task",
			taskID: 1,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					UpdateTask(ctx, 1).
					Return(nil)
			},
			expectedResponse: nil,
		},
		{
			name:   "Error - Database failure",
			taskID: 1,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					UpdateTask(ctx, 1).
					Return(errors.New("database error"))
			},
			expectedResponse: errors.New("database error"),
		},
		{
			name:   "Error - Task not found",
			taskID: 999,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					UpdateTask(ctx, 999).
					Return(errors.New("task not found"))
			},
			expectedResponse: errors.New("task not found"),
		},
		{
			name:   "Error - Invalid task ID",
			taskID: 0,
			mockExpect: func(ctx *gofr.Context) {
				// No call to store expected
			},
			expectedResponse: errors.New("invalid task ID"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &gofr.Context{
				Context: context.Background(),
			}

			if tt.mockExpect != nil {
				tt.mockExpect(ctx)
			}

			s := NewTaskService(mockStore)

			err := s.UpdateTask(ctx, tt.taskID)

			assert.Equal(t, tt.expectedResponse, err, "TEST[%d] - %s", i, tt.name)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockTaskStore(ctrl)

	tests := []struct {
		name             string
		taskID           int
		mockExpect       func(ctx *gofr.Context)
		expectedResponse error
	}{
		{
			name:   "Success - Delete task",
			taskID: 1,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					DeleteTask(ctx, 1).
					Return(nil)
			},
			expectedResponse: nil,
		},
		{
			name:   "Error - Database failure",
			taskID: 1,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					DeleteTask(ctx, 1).
					Return(errors.New("database error"))
			},
			expectedResponse: errors.New("database error"),
		},
		{
			name:   "Error - Task not found",
			taskID: 999,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					DeleteTask(ctx, 999).
					Return(errors.New("task not found"))
			},
			expectedResponse: errors.New("task not found"),
		},
		{
			name:   "Error - Invalid task ID",
			taskID: 0,
			mockExpect: func(ctx *gofr.Context) {
				// No store call for validation errors
			},
			expectedResponse: errors.New("invalid task ID"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock Gofr context
			ctx := &gofr.Context{
				Context: context.Background(),
			}

			if tt.mockExpect != nil {
				tt.mockExpect(ctx)
			}

			s := NewTaskService(mockStore)

			err := s.DeleteTask(ctx, tt.taskID)

			assert.Equal(t, tt.expectedResponse, err, "TEST[%d] - %s", i, tt.name)
		})
	}
}
