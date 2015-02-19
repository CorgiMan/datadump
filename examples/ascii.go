package main

import "github.com/CorgiMan/datadump"

func main() {
	datadump.Open(":8080")
	defer datadump.Close()
	mat := []interface{}{"xxxx", "xxox", "xxxx"}
	for i := range mat {
		mat[i] = []byte(mat[i].(string))
	}
	datadump.Show(mat)
}
