package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/saurabh-sde/employee-go/model"
	"github.com/saurabh-sde/employee-go/utility"
	"github.com/spf13/cast"
)

var mutex sync.Mutex

// CreateEmployee handler function
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	// For simplicity, just use a counter on len of map
	emp.ID = len(model.EmployeesDB) + 1

	model.EmployeesDB[emp.ID] = emp

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

// GetEmployeeByID handler function
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		utility.Print("invalid employee id")
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	emp, ok := model.EmployeesDB[cast.ToInt(id)]
	if !ok {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(emp)
}

// UpdateEmployee handler function
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := model.EmployeesDB[emp.ID]; !ok {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	model.EmployeesDB[emp.ID] = emp

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

// DeleteEmployee handler function
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	empID, ok := mux.Vars(r)["id"]
	if !ok {
		utility.Print("invalid employee id")
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	_, ok = model.EmployeesDB[cast.ToInt(empID)]
	if !ok {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	delete(model.EmployeesDB, cast.ToInt(empID))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted Employee: " + empID)
}

func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	// get query params for pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // default page size
	}

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	// prepare employees arr
	employees := []model.Employee{}
	for _, emp := range model.EmployeesDB {
		employees = append(employees, emp)

	}
	// returning resp based employee ID between range is wrong as there is delete emp also
	// so prepare arr of emp sort it on empId and then return required resp b/w resp = employees[start:end]

	// sort according to empID
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].ID < employees[j].ID
	})

	if startIndex < len(employees) {
		if endIndex >= len(employees) {
			// take rest of entries
			employees = employees[startIndex:]
		} else {
			employees = employees[startIndex:endIndex]
		}
	} else {
		// start out of bounds so return blank resp
		employees = []model.Employee{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employees)
}

// Gin Router Handlers Functions

// CreateEmployee handler function
func CreateEmployeeGin(c *gin.Context) {
	var emp model.Employee
	err := c.Bind(&emp)
	if err != nil {
		utility.Print("CreateEmployeeGin: Error in request: ", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	// For simplicity, just use a counter on len of map
	emp.ID = len(model.EmployeesDB) + 1

	model.EmployeesDB[emp.ID] = emp
	c.JSON(http.StatusOK, gin.H{"employee": emp})
}

// GetEmployeeByID handler function
func GetEmployeeByIDGin(c *gin.Context) {
	id := c.Param("id")

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	emp, ok := model.EmployeesDB[cast.ToInt(id)]
	if !ok {
		err := errors.New("Employee not found")
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, emp)
}

// UpdateEmployee handler function
func UpdateEmployeeGin(c *gin.Context) {
	var emp model.Employee
	err := c.Bind(&emp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := model.EmployeesDB[emp.ID]; !ok {
		err = errors.New("Employee not found")
		c.JSON(http.StatusNotFound, err)
		return
	}

	model.EmployeesDB[emp.ID] = emp

	c.JSON(http.StatusOK, emp)
}

// DeleteEmployee handler function
func DeleteEmployeeGin(c *gin.Context) {
	empID := c.Param("id")

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := model.EmployeesDB[cast.ToInt(empID)]
	if !ok {
		err := errors.New("Employee not found")
		c.JSON(http.StatusNotFound, err)
		return
	}

	delete(model.EmployeesDB, cast.ToInt(empID))

	c.JSON(http.StatusOK, "Deleted Employee: "+empID)
}

func GetAllEmployeesGin(c *gin.Context) {
	// get query params for pagination
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // default page size
	}

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	// acquire lock to handle concurrent race conditions
	mutex.Lock()
	defer mutex.Unlock()

	// prepare employees arr
	employees := []model.Employee{}
	for _, emp := range model.EmployeesDB {
		employees = append(employees, emp)

	}
	// returning resp based employee ID between range is wrong as there is delete emp also
	// so prepare arr of emp sort it on empId and then return required resp b/w resp = employees[start:end]

	// sort according to empID
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].ID < employees[j].ID
	})

	if startIndex < len(employees) {
		if endIndex >= len(employees) {
			// take rest of entries
			employees = employees[startIndex:]
		} else {
			employees = employees[startIndex:endIndex]
		}
	} else {
		// start out of bounds so return blank resp
		employees = []model.Employee{}
	}

	c.JSON(http.StatusOK, gin.H{
		"employees": employees,
	})
}
