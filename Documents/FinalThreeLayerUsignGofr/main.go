package main

import (
	_ "FinalThreeLayerUsingGOfr/docs"
	taskHandler2 "FinalThreeLayerUsingGOfr/handler/taskHandler"
	"FinalThreeLayerUsingGOfr/handler/userHandler"
	taskStore "FinalThreeLayerUsingGOfr/store/task"
	userStore "FinalThreeLayerUsingGOfr/store/user"
	_ "github.com/swaggo/http-swagger"
	// httpSwagger "github.com/swaggo/http-swagger"
	"gofr.dev/pkg/gofr"

	userService "FinalThreeLayerUsingGOfr/service"
	//"log"
	//"net/http"
)

func main() {
	// db := dataSource.InitDB("root:7280@tcp(localhost:3306)/taskService")

	app := gofr.New()

	tStore := taskStore.New()
	uStore := userStore.New()

	tService := userService.NewTaskService(tStore)
	uService := userService.NewUserService(uStore)

	tHandler := taskHandler2.NewTaskHandler(tService)
	uHandler := userHandler.NewUserHandler(uService)

	//http.Handle("/swagger/", httpSwagger.WrapHandler)
	//http.HandleFunc("GET /tasks", tHandler.GetAll)
	//http.HandleFunc("GET /tasks/{id}", tHandler.GetById)
	//http.HandleFunc("POST /tasks", tHandler.CreateTask)
	//http.HandleFunc("DELETE /tasks/{id}", tHandler.DeleteById)
	//http.HandleFunc("PUT /tasks/{id}", tHandler.MarkCompleteById)
	//http.HandleFunc("POST /users", uHandler.CreateUser)
	//http.HandleFunc("GET /users/{id}", uHandler.GetByUserId)
	//
	//log.Println("services running on :8081")
	//log.Fatal(http.ListenAndServe(":8081", nil))

	app.GET("/tasks", tHandler.GetAll)
	app.GET("/tasks/{id}", tHandler.GetByID)
	app.POST("/tasks", tHandler.Create)
	app.DELETE("/tasks/{id}", tHandler.Delete)
	app.PUT("/tasks/{id}", tHandler.MarkCompleteById)
	app.POST("/user", uHandler.Create)
	app.GET("/user/{id}", uHandler.GetByID)

	app.Run()

}
