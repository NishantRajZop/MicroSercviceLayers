package userHandler

import (
	_ "FinalThreeLayerUsingGOfr/docs" // Import docs package for Swagger
	"FinalThreeLayerUsingGOfr/models"
	"errors"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"strconv"
)

type userHandler struct {
	service userService
}

func NewUserHandler(service userService) *userHandler {
	return &userHandler{service: service}
}

// @Summary Create a new user
// @Description Add a new user to the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 "User created successfully"
// @Failure 400 {object} object "Invalid input/Empty fields"
// @Failure 500 {object} object "Internal Server Error"
// @Router /users [post]
func (h *userHandler) Create(ctx *gofr.Context) (interface{}, error) {
	var user models.User

	// Bind request body
	if err := ctx.Bind(&user); err != nil {
		return nil, gofrHttp.ErrorInvalidParam{
			Params: []string{"Invalid Json"}} // Returns HTTP 400
	}

	// Validate required fields
	if user.Name == "" {
		return nil, errors.New("name is required")
	}

	// Business validation (example: ID should not be provided for creation)
	//if user.ID != 0 {
	//	return nil, gofrHttp.ErrorInvalidParam{
	//		Params: strings.Split("id", "&"),
	//	}
	//}

	// Call service layer

	var resUser models.User
	resUser, err := h.service.CreateUser(ctx, user)
	if err != nil {
		return nil, gofrHttp.ErrorEntityAlreadyExist{}
	}
	// Return created user with HTTP 201
	return resUser, nil
}

// @Summary Get user by ID
// @Description Get user details by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User "User details"
// @Failure 400 {object} object "Invalid ID format"
// @Failure 404 {object} object "User not found"
// @Router /users/{id} [get]
func (h *userHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	// Get and validate ID from path parameter
	reqId := ctx.PathParam("id")
	id, err := strconv.Atoi(reqId)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}} // Returns HTTP 400
	}
	// Get user from service layer
	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		return nil, gofrHttp.ErrorEntityNotFound{
			Name:  "user",
			Value: strconv.Itoa(id),
		} // Returns HTTP 404
	}
	// Returns HTTP 200 with user JSON
	return user, nil
}
