package main

import (
	"os"

	"github.com/CorgiMan/datadump"
)

func main() {
	datadump.Open(":8080")
	defer datadump.Close()

	// // print a string
	datadump.Show("Hello World!")

	// // print a type (a file for example)
	f, _ := os.Open("main.go")
	datadump.Show(f)
}
