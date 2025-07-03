package userHandler

import (
	"FinalThreeLayerUsingGOfr/models"
	"bytes"
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

func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockuserService(ctrl)
	handler := NewUserHandler(mockSvc)

	tests := []struct {
		name             string
		requestBody      string
		mockExpect       func(ctx *gofr.Context)
		expectedStatus   int
		expectedResponse interface{}
		expectedError    error
	}{
		{
			name:        "Success - Create user",
			requestBody: `{"name":"John Doe"}`,
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					CreateUser(ctx, models.User{Name: "John Doe"}).
					Return(models.User{Name: "John Doe"}, nil)
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: models.User{ID: 0, Name: "John Doe"},
		},
		{
			name:           "Error - Invalid payload",
			requestBody:    `{invalid json}`,
			mockExpect:     func(ctx *gofr.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  gofrHttp.ErrorInvalidParam{Params: []string{"Invalid Json"}},
		},
		{
			name:           "Error - Missing name",
			requestBody:    `{"id":1}`,
			mockExpect:     func(ctx *gofr.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.New("name is required"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			ctx := &gofr.Context{
				Context: context.Background(),
				Request: gofrHttp.NewRequest(req),
			}

			tt.mockExpect(ctx)

			resp, err := handler.Create(ctx)

			if tt.expectedResponse != nil {
				assert.Equal(t, tt.expectedResponse, resp, "TEST[%d] - %s", i, tt.name)
			}

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err, "TEST[%d] - %s", i, tt.name)
			}
		})
	}
}

func TestUserHandler_GetByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockuserService(ctrl)
	handler := NewUserHandler(mockSvc)

	tests := []struct {
		name             string
		userID           string
		mockExpect       func(ctx *gofr.Context)
		expectedStatus   int
		expectedResponse interface{}
		expectedError    error
	}{
		{
			name:   "Success - Get user by ID",
			userID: "2",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					GetUserByID(ctx, 2).
					Return(models.User{ID: 2, Name: "Nishant"}, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: models.User{ID: 2, Name: "Nishant"},
		},
		{
			name:           "Error - Invalid user ID",
			userID:         "abc",
			mockExpect:     func(ctx *gofr.Context) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:   "Error - User not found",
			userID: "999",
			mockExpect: func(ctx *gofr.Context) {
				mockSvc.EXPECT().
					GetUserByID(ctx, 999).
					Return(models.User{}, gofrHttp.ErrorEntityNotFound{Name: "user", Value: "999"})
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  gofrHttp.ErrorEntityNotFound{Name: "user", Value: "999"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": tt.userID,
			})

			ctx := &gofr.Context{
				Context: context.Background(),
				Request: gofrHttp.NewRequest(req),
			}

			tt.mockExpect(ctx)

			_, err := handler.GetByID(ctx)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err, "TEST[%d] - %s", i, tt.name)
			}
		})
	}
}
