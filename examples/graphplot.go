package main

import (
	"math"

	"github.com/CorgiMan/datadump"
)

func main() {
	datadump.Open(":8080")
	defer datadump.Close()

	// plot a sine function
	xs := []float64{}
	ys := []float64{}
	for x := -30.0; x < 30.0; x += 0.1 {
		xs = append(xs, x)
		ys = append(ys, math.Sin(x))
	}
	datadump.Show(map[string]interface{}{"connected": 0, "x": xs, "y": ys})
	xs = []float64{}
	ys = []float64{}
	for t := 0.0; t < 200; t += 0.01 {
		xs = append(xs, math.Sin(t*1.17))
		ys = append(ys, math.Sin(t))
	}
	datadump.Show(map[string]interface{}{"connected": 0, "x": xs, "y": ys})
}
