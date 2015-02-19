package main

// show a [][]number as a grayscale or show an image.Image with datadump

import (
	"image"
	"image/color"
	"math/cmplx"

	"github.com/CorgiMan/datadump"
)

func main() {
	datadump.Open(":8080")
	defer datadump.Close()

	matrix := make([][]int, 16)
	for i := range matrix {
		matrix[i] = make([]int, 16)
		for j := range matrix[i] {
			matrix[i][j] = i ^ j
		}
	}
	datadump.Show(matrix)

	datadump.Show(Mandelbrot{Width: 300, Height: 300})
}

// implements image.Image
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
			col := uint8((i * 100) % 255)
			return color.RGBA{0, col, col, 255}
		}
	}
	return color.RGBA{0, 0, 0, 255}
}

func (img Mandelbrot) ColorModel() color.Model {
	return color.RGBAModel
}
