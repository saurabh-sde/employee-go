package model

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/saurabh-sde/employee-go/utility"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var EmployeesDB map[int]Employee
var Db *mongo.Database

type Employee struct {
	// ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EmployeeID int     `bson:"employee_id,omitempty" json:"employeeId,omitempty"` // omitempty to handle create emp
	Name       string  `bson:"name" json:"name"`
	Position   string  `bson:"position" json:"position"`
	Salary     float64 `bson:"salary" json:"salary"`
}

func init() {
	utility.Print("Connecting DB")
	EmployeesDB = make(map[int]Employee)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading configs")
		return
	}
	uri := os.Getenv("DB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// database and colletion code goes here
	db := client.Database("dev")

	Db = db
	// defer func() {
	// 	utility.Print("Disconnecting DB")
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
}
