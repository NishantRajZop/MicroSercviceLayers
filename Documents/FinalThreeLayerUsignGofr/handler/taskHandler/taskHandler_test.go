package taskHandler

import (
	"FinalThreeLayerUsingGOfr/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMocktaskService(ctrl)
	handler := NewTaskHandler(mockSvc)

	expectedTasks := []models.Task{
		{ID: 1, Title: "Task 1", UserID: 1, Completed: false},
		{ID: 2, Title: "Task 2", UserID: 1, Completed: true},
	}

	tests := []struct {
		name             string
		mockExpect       func(ctx *gofr.Context)
		expectedStatus   int
		expectedResponse interface{}
		expectedError    error
	}{
		{
			name: "Success - Get all tasks",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					GetAllTasks(ctx).
					Return(expectedTasks, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: expectedTasks,
		},
		{
			name: "Error - Service failure",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					GetAllTasks(ctx).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  errors.New("database error"),
		},
		{
			name: "Success - Empty result",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					GetAllTasks(ctx).
					Return([]models.Task{}, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: []models.Task{},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			//w := httptest.NewRecorder()

			// Create GoFr context
			ctx := &gofr.Context{
				Context: context.Background(),
				Request: gofrHttp.NewRequest(req),
			}

			// Set expectations on the service
			tt.mockExpect(ctx)

			// Call handler
			resp, err := handler.GetAll(ctx)

			// Response assertions
			if tt.expectedResponse != nil {
				assert.Equal(t, tt.expectedResponse, resp, "TEST[%d], Failed.\n%s", i, tt.name)
			}
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error(), "TEST[%d], Failed.\n%s", i, tt.name)
			} else {
				assert.NoError(t, err, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}

func TestTaskHandler_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMocktaskService(ctrl)
	handler := NewTaskHandler(mockSvc)

	expectedTask := models.Task{
		ID:        1,
		Title:     "Task 1",
		UserID:    1,
		Completed: false,
	}

	tests := []struct {
		name             string
		taskID           string
		mockExpect       func(ctx *gofr.Context)
		expectedStatus   int
		expectedResponse interface{}
		expectedError    error
	}{
		{
			name:   "Success - Get task by ID",
			taskID: "1",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					GetTaskByID(ctx, 1).
					Return(expectedTask, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: expectedTask,
		},
		{
			name:           "Error - Invalid ID",
			taskID:         "abc",
			mockExpect:     func(ctx *gofr.Context) {}, // No call expected
			expectedStatus: http.StatusBadRequest,
			expectedError:  gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:   "Error - Task not found",
			taskID: "999",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					GetTaskByID(ctx, 999).
					Return(models.Task{}, sql.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  gofrHttp.ErrorEntityNotFound{Name: "taskId", Value: "999"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.taskID, nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.taskID,
			})

			// Create GoFr context
			ctx := &gofr.Context{
				Context: context.Background(),
				Request: gofrHttp.NewRequest(req),
			}

			// Run mock expectation
			tt.mockExpect(ctx)

			// Call handler
			resp, err := handler.GetByID(ctx)

			// Assertions
			if tt.expectedResponse != nil {
				assert.Equal(t, tt.expectedResponse, resp, "TEST[%d], Failed.\n%s", i, tt.name)
			}
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err, "TEST[%d], Failed.\n%s", i, tt.name)
			} else {
				assert.NoError(t, err, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}

func TestTaskHandler_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMocktaskService(ctrl)
	handler := NewTaskHandler(mockSvc)

	tests := []struct {
		name             string
		requestBody      string
		mockExpect       func(ctx *gofr.Context)
		expectedStatus   int
		expectedResponse string
		expectedError    error
	}{
		{
			name:        "Success - Create task",
			requestBody: `{ "id":1 ,"title" : "DemoTask 1", "user_id" : 101}`,
			mockExpect: func(ctx *gofr.Context) {
				task := models.Task{ID: 1, Title: "DemoTask 1", UserID: 101}
				mockSvc.EXPECT().
					CreateTask(ctx, task).
					Return(task, nil)
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: `{"id": 1,"title": "DemoTask 1","user_id": 101,"completed": false}`,
		},
		{
			name:           "Error - Invalid payload",
			requestBody:    `{invalid json}`,
			mockExpect:     func(ctx *gofr.Context) {}, // No call expected
			expectedStatus: http.StatusBadRequest,
			expectedError:  gofrHttp.ErrorInvalidParam{Params: []string{"body", "query"}},
		},
		{
			name:           "Error - Missing title",
			requestBody:    `{"description":"Task desc"}`,
			mockExpect:     func(ctx *gofr.Context) {}, // No Call expected here too becasue in Handler layer only if task,Title is empty we retrun
			expectedStatus: http.StatusNotFound,
			expectedError:  gofrHttp.ErrorEntityNotFound{Name: "title"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create GoFr context
			ctx := &gofr.Context{
				Context: context.Background(),
				Request: gofrHttp.NewRequest(req),
			}

			// Run mock expectation
			tt.mockExpect(ctx)

			// Call the handler
			resp, err := handler.Create(ctx)

			// Assertions
			if tt.expectedResponse != "" {
				responseBody, _ := json.Marshal(resp)
				assert.JSONEq(t, tt.expectedResponse, string(responseBody), "TEST[%d], Failed.\n%s", i, tt.name)
			}
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err, "TEST[%d], Failed.\n%s", i, tt.name)
			} else {
				assert.NoError(t, err, "TEST[%d], Failed.\n%s", i, tt.name)
			}
		})
	}
}

func TestTaskHandler_MarkCompleteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMocktaskService(ctrl)
	handler := NewTaskHandler(mockSvc)

	tests := []struct {
		name           string
		taskID         string
		mockExpect     func(ctx *gofr.Context)
		expectedStatus int
		expectedError  error
	}{
		{
			name:   "Success - Mark task complete",
			taskID: "8",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					UpdateTask(ctx, 8).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Invalid ID",
			taskID:         "abc",
			mockExpect:     func(ctx *gofr.Context) {}, // no service call
			expectedStatus: http.StatusBadRequest,
			expectedError:  gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:   "Error - Task not found",
			taskID: "999",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					UpdateTask(gomock.Any(), 999).
					Return(gofrHttp.ErrorEntityNotFound{Name: "Either The id is Already Marked or Not Found", Value: "999"})
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  gofrHttp.ErrorEntityNotFound{Name: "Either The id is Already Marked or Not Found", Value: "999"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.taskID,
			})
			ctx := &gofr.Context{
				Context: context.Background(),
				Request: gofrHttp.NewRequest(req),
			}

			tt.mockExpect(ctx)

			// Call handler
			_, err := handler.MarkCompleteById(ctx)

			// Assert error
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err, "TEST[%d] - %s", i, tt.name)
			} else {
				assert.NoError(t, err, "TEST[%d] - %s", i, tt.name)
			}
		})
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMocktaskService(ctrl)
	handler := NewTaskHandler(mockSvc)

	tests := []struct {
		name           string
		taskID         string
		mockExpect     func(ctx *gofr.Context)
		expectedStatus int
		expectedError  error
	}{
		{
			name:   "Success - Delete task",
			taskID: "8",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					DeleteTask(ctx, 8).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
		{
			name:   "Error - Invalid ID",
			taskID: "abc",
			mockExpect: func(ctx *gofr.Context) {
				// No service call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:   "Error - Task not found",
			taskID: "999",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					DeleteTask(gomock.Any(), 999).
					Return(gofrHttp.ErrorEntityNotFound{Name: "task", Value: "999"})
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  gofrHttp.ErrorEntityNotFound{Name: "task", Value: "999"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tt.taskID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.taskID})

			ctx := &gofr.Context{
				Context: context.Background(),
				Request: gofrHttp.NewRequest(req),
			}

			tt.mockExpect(ctx)

			// Call the handler
			_, err := handler.Delete(ctx)

			// Assert error
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err, "TEST[%d] - %s", i, tt.name)
			} else {
				assert.NoError(t, err, "TEST[%d] - %s", i, tt.name)
			}
		})
	}
}
