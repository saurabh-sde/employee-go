package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saurabh-sde/employee-go/model"
)

func TestCreateEmployee(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name        string
		args        args
		expectedID  int
		expectedErr bool
	}{
		// Added test cases.
		{
			name: "Valid Request",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/create", bytes.NewBuffer([]byte(`{"Name":"John","Position":"Software Engineer","Salary":60000}`))),
			},
			expectedID:  1,
			expectedErr: false,
		},
		{
			name: "Invalid Request-Empty Body",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/create", nil),
			},
			expectedID:  0,
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateEmployee(tt.args.w, tt.args.r)

			// Check the response status code
			if rr, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if tt.expectedErr && rr.Code != http.StatusBadRequest {
					t.Errorf("Expected status code %d but got %d", http.StatusBadRequest, rr.Code)
				}
				if !tt.expectedErr && rr.Code != http.StatusOK {
					t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
				}
			} else {
				t.Error("ResponseRecorder not available")
			}

			// If the request was successful, check if the employee was created with the expected ID
			if !tt.expectedErr {
				if rr, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
					var emp model.Employee
					err := json.NewDecoder(rr.Body).Decode(&emp)
					if err != nil {
						t.Errorf("Error decoding response body: %v", err)
					}
					fmt.Println("added employee: ", emp)
					if emp.ID != tt.expectedID {
						t.Errorf("Expected employee ID %d but got %d", tt.expectedID, emp.ID)
					}
				} else {
					t.Error("ResponseRecorder not available")
				}
			}
		})
	}
}
