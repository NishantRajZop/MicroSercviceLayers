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

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserStore(ctrl)
	s := NewUserService(mockStore)

	tests := []struct {
		name             string
		user             models.User
		mockExpect       func(ctx *gofr.Context)
		expectedResult   models.User
		expectedResponse error
	}{
		{
			name: "Success - User created",
			user: models.User{ID: 1, Name: "John Doe"}, // Don't set ID here; it’s generated
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					CreateUser(ctx, models.User{ID: 1, Name: "John Doe"}).
					Return(models.User{ID: 1, Name: "John Doe"}, nil)
			},
			expectedResult:   models.User{ID: 1, Name: "John Doe"},
			expectedResponse: nil,
		},
		{
			name:             "Error - Empty name",
			user:             models.User{Name: ""},
			mockExpect:       nil, // Validation error: store not called
			expectedResult:   models.User{},
			expectedResponse: errors.New("user name cannot be empty"),
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

			_, err := s.CreateUser(ctx, tt.user)

			assert.Equal(t, tt.expectedResponse, err, "TEST[%d] - %s", i, tt.name)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockUserStore(ctrl)
	s := NewUserService(mockStore)

	tests := []struct {
		name             string
		userID           int
		mockExpect       func(ctx *gofr.Context)
		expectedResponse models.User
		expectedError    error
	}{
		{
			name:   "Success - User found",
			userID: 1,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					GetUserByID(ctx, 1).
					Return(models.User{ID: 1, Name: "Nishant"}, nil)
			},
			expectedResponse: models.User{ID: 1, Name: "Nishant"},
			expectedError:    nil,
		},
		{
			name:   "Error - User not found",
			userID: 999,
			mockExpect: func(ctx *gofr.Context) {
				mockStore.EXPECT().
					GetUserByID(ctx, 999).
					Return(models.User{}, sql.ErrNoRows)
			},
			expectedResponse: models.User{},
			expectedError:    sql.ErrNoRows,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &gofr.Context{
				Context: context.Background(),
			}

			tt.mockExpect(ctx)

			user, err := s.GetUserByID(ctx, tt.userID)

			assert.Equal(t, tt.expectedResponse, user, "TEST[%d] - %s", i, tt.name)
			assert.Equal(t, tt.expectedError, err, "TEST[%d] - %s", i, tt.name)
		})
	}
}
