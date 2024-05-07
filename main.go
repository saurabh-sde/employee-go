package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/saurabh-sde/employee-go/handler"
	"github.com/saurabh-sde/employee-go/middleware"
)

func init() {

}

func main() {
	r := mux.NewRouter()

	// employee route
	r.HandleFunc("/employee", handler.CreateEmployee).Methods("POST")
	r.HandleFunc("/employee", handler.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employee/{id}", handler.GetEmployeeByID).Methods("Get")
	r.HandleFunc("/employee/{id}", handler.DeleteEmployee).Methods("Delete")
	r.HandleFunc("/employees", handler.GetAllEmployees).Methods("Get")

	// add loggin middleware
	r.Use(middleware.Logging)

	fmt.Println("Starting Local Server: http://localhost:8080")
	// start local server
	http.ListenAndServe(":8080", r)
}
