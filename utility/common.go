package utility

import (
	"fmt"
)

func Print(v ...interface{}) {
	logLvl := "INFO"
	fmt.Printf("%s: %+v\n", logLvl, v)
}

func Error(v ...interface{}) {
	logLvl := "ERROR"
	fmt.Printf("%s: %+v\n", logLvl, v)
}
