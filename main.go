package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/saurabh-sde/employee-go/handler"
	"github.com/saurabh-sde/employee-go/middleware"
)

func init() {

}

var listen chan struct{}

func main() {

	listen = make(chan struct{})

	r := mux.NewRouter()
	// employee route
	r.HandleFunc("/employee", handler.CreateEmployee).Methods("POST")
	r.HandleFunc("/employee", handler.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employee/{id}", handler.GetEmployeeByID).Methods("Get")
	r.HandleFunc("/employee/{id}", handler.DeleteEmployee).Methods("Delete")
	r.HandleFunc("/employees", handler.GetAllEmployees).Methods("Get")
	// add loggin middleware
	r.Use(middleware.Logging)
	// start local server
	srv1 := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r2 := gin.Default()
	r2.POST("/employee", handler.CreateEmployeeGin)
	r2.PUT("/employee", handler.UpdateEmployeeGin)
	r2.GET("/employee/:id", handler.GetEmployeeByIDGin)
	r2.DELETE("/employee/:id", handler.DeleteEmployeeGin)
	r2.GET("/employees", handler.GetAllEmployeesGin)

	srv2 := &http.Server{
		Addr:    ":8000",
		Handler: r2,
	}

	go func() {
		fmt.Println("Starting Local Server: http://localhost:8080")
		if err := srv1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe() error: %v\n", err)
		}
	}()

	go func() {
		fmt.Println("Starting Local Server: http://localhost:8000")
		if err := srv2.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe() error: %v\n", err)
		}
	}()

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Channel to signal the completion of shutdown
	done := make(chan struct{})

	// Goroutine to handle shutdown
	go func() {
		<-stop

		// Create a deadline to wait for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv1.Shutdown(ctx); err != nil {
			fmt.Printf("Server 1 Shutdown Failed:%+v", err)
		}
		if err := srv2.Shutdown(ctx); err != nil {
			fmt.Printf("Server 2 Shutdown Failed:%+v", err)
		}

		// Signal that shutdown is complete
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Servers shut down gracefully")
	case <-stop:
		fmt.Println("Received second signal, shutting down forcefully")
	}

	// Close the listen channel to signal the end of the application
	close(listen)

	fmt.Println("Shutting down gracefully, press Ctrl+C again to force")
}
