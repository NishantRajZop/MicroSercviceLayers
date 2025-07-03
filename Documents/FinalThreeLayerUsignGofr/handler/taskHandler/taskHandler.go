package taskHandler

import (
	"fmt"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"strconv"

	_ "FinalThreeLayerUsingGOfr/docs" // Import docs package for Swagger
	"FinalThreeLayerUsingGOfr/models"
)

type taskHandler struct {
	service taskService
}

func NewTaskHandler(service taskService) *taskHandler {
	return &taskHandler{service: service}
}

// @Summary Get all tasks
// @Description Retrieve a list of all tasks
// @Tags Tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.Task "List of tasks"
// @Failure 500 {object} object "Internal Server Error"
// @Router /tasks [get]
// GetAll returns all tasks with proper error handling and response formatting
func (h *taskHandler) GetAll(ctx *gofr.Context) (interface{}, error) {
	// Get tasks from service layer
	tasks, err := h.service.GetAllTasks(ctx)

	// Error handling - automatically sets appropriate status code
	if err != nil {
		// Return error - GoFr will automatically:
		// 1. Set Content-Type: application/json
		// 2. Set HTTP 500 status code
		// 3. Format error response consistently
		return nil, err
	}

	// Successful response - GoFr automatically:
	// 1. Sets Content-Type: application/json
	// 2. Sets HTTP 200 status code
	// 3. Properly encodes the response
	return tasks, nil
}

// @Summary Get a task by ID
// @Description Get a single task by its ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task "Task details"
// @Failure 400 {object} object "Invalid ID format"
// @Failure 404 {object} object "Task not found"
// @Router /tasks/{id} [get]
func (h *taskHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	// Get and validate ID from path parameter
	reqId := ctx.Request.PathParam("id")
	id, err := strconv.Atoi(reqId)
	if err != nil {
		// Correct way to return invalid parameter error
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}
	// Get task from service layer
	task, err := h.service.GetTaskByID(ctx, id)
	if err != nil {
		// Returns HTTP 404 with {"error":"task not found","entity":"task","id":"123"}
		return nil, gofrHttp.ErrorEntityNotFound{Name: "taskId", Value: reqId}
	}

	// Returns HTTP 200 with the task JSON
	return task, nil
}

// @Summary Create a new task
// @Description Add a new task to the system
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task data"
// @Success 201 {object} models.Task "Task created successfully"
// @Failure 400 {object} object "Invalid input/Empty title"
// @Failure 500 {object} object "Internal Server Error"
// @Router /tasks [post]
func (h *taskHandler) Create(ctx *gofr.Context) (interface{}, error) {
	var task models.Task

	// Bind request body to task struct
	if err := ctx.Bind(&task); err != nil {
		fmt.Println("Bind error:", err)
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"body", "query"}}
	}
	fmt.Printf("%+v\n", task)

	// Validate required fields
	if task.Title == "" {
		return nil, gofrHttp.ErrorEntityNotFound{Name: "title"}
	}

	// Call service layer
	if _, err := h.service.CreateTask(ctx, task); err != nil {
		return nil, err // Will return 500 by default for generic errors
	}

	// Return created task with 201 status
	return task, nil
}

// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204 "Task deleted successfully"
// @Failure 400 {object} object "Invalid ID format"
// @Failure 500 {object} object "Internal Server Error"
// @Router /tasks/{id} [delete]

func (h *taskHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	// Get and validate ID from path parameter
	reqId := ctx.PathParam("id")
	id, err := strconv.Atoi(reqId)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}} // Returns HTTP 400
	}

	// Call service layer
	if err := h.service.DeleteTask(ctx, id); err != nil {
		return nil, gofrHttp.ErrorEntityNotFound{
			Name:  "task",
			Value: strconv.Itoa(id),
		} // Returns HTTP 404
	}

	// Return nil for successful deletion (GoFr will send HTTP 204)
	return nil, nil
}

//	@Summary Mark task as complete
//
// @Description Update a task's status to complete
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 "Task marked as complete"
// @Failure 400 {object} object "Invalid ID format"
// @Failure 500 {object} object "Internal Server Error"
// @Router /tasks/{id} [put]
func (h *taskHandler) MarkCompleteById(ctx *gofr.Context) (interface{}, error) {
	reqId := ctx.PathParam("id")
	fmt.Println("reqId", reqId)
	id, err := strconv.Atoi(reqId)
	fmt.Println("id", id)
	fmt.Println("err", err)
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}} // Returns HTTP 400
	}

	// Call service layer to mark task complete
	err = h.service.UpdateTask(ctx, id)
	if err != nil {
		return nil, gofrHttp.ErrorEntityNotFound{
			Name:  "Either The id is Already Marked or Not Found",
			Value: strconv.Itoa(id),
		} // Returns HTTP 404
	}
	// Return the updated task with HTTP 200
	return models.Task{
		ID:        id,
		Completed: true,
	}, nil
}
