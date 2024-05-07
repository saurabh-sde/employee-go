package model

var EmployeesDB map[int]Employee

type Employee struct {
	ID       int     `json:"id,omitempty"` // omitempty to handle create emp
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

func init() {
	EmployeesDB = make(map[int]Employee)
}
