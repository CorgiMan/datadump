package main

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"
	"net/http"
	"os"

	"github.com/CorgiMan/datadump"
)

type Mandelbrot struct {
	Width, Height int
}

func (img Mandelbrot) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

func (img Mandelbrot) At(x, y int) color.Color {
	c := complex(float64(x)/float64(img.Width)*3.0-2.0, float64(y)/float64(img.Height)*3.0-1.5)
	z := complex(0, 0)
	i := 0
	for ; i < 1000; i++ {
		z = z*z + c
		if cmplx.Abs(z) > 2 {
			return color.Gray{uint8((i * 100) % 255)}
		}
	}
	return color.Gray{0}
}

func (img Mandelbrot) ColorModel() color.Model {
	return color.GrayModel
}

func main() {
	datadump.Open()
	defer datadump.Close()
	datadump.C <- Mandelbrot{100, 100}

	ys := []float64{}
	xs := []float64{}
	for i := -300; i < 300; i++ {
		xs = append(xs, float64(i)*0.1)
		ys = append(ys, math.Sin(float64(xs[len(xs)-1]))*6)
	}
	datadump.C <- map[string]interface{}{"connected": 0, "x": xs, "y": ys}
	res, err := http.Get("http://www.asterank.com/api/skymorph/search?target=J99TS7A")
	jsonquery
	datadump.C <- res
	datadump.C <- err
	r, err := os.Open("main.go")
	datadump.C <- r
	datadump.C <- err

}
