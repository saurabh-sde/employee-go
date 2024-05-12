package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/saurabh-sde/employee-go/handler"
)

func init() {

}

func main() {
	// r := mux.NewRouter()
	// // employee route
	// r.HandleFunc("/employee", handler.CreateEmployee).Methods("POST")
	// r.HandleFunc("/employee", handler.UpdateEmployee).Methods("PUT")
	// r.HandleFunc("/employee/{id}", handler.GetEmployeeByID).Methods("Get")
	// r.HandleFunc("/employee/{id}", handler.DeleteEmployee).Methods("Delete")
	// r.HandleFunc("/employees", handler.GetAllEmployees).Methods("Get")
	// // add loggin middleware
	// r.Use(middleware.Logging)
	// start local server
	// http.ListenAndServe(":8080", r)

	r := gin.Default()

	r.POST("/employee", handler.CreateEmployeeGin)
	r.PUT("/employee", handler.UpdateEmployeeGin)
	r.GET("/employee/:id", handler.GetEmployeeByIDGin)
	r.DELETE("/employee/:id", handler.DeleteEmployeeGin)
	r.GET("/employees", handler.GetAllEmployeesGin)

	fmt.Println("Starting Local Server: http://localhost:8000")
	r.Run(":8000")
}
