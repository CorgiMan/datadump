package datadump

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
)

func write2DImage(w http.ResponseWriter, v interface{}) (ok bool, err error) {
	a, ok := float2D(v)
	if !ok {
		return
	}

	if len(a) == 0 || len(a[0]) == 0 {
		return
	}
	min, max := a[0][0], a[0][0]
	for i := range a {
		for j := range a[i] {
			if a[i][j] > max {
				max = a[i][j]
			}
			if a[i][j] < min {
				min = a[i][j]
			}
		}
	}
	pixs := make([][]uint8, len(a))
	for i := range a {
		pixs[i] = make([]uint8, len(a[i]))
		for j := range a[i] {
			pixs[i][j] = uint8(255 * (a[i][j] - min) / (max - min))
		}
	}

	img := StretchGrayImage{pixs, 350 * len(a[0]) / len(a), 350}

	return true, writeHtmlImg(w, img)
}

func float2D(v interface{}) ([][]float64, bool) {
	rows, ok := v.([]interface{})
	fmt.Println(ok)
	if !ok {
		return nil, false
	}

	a := make([][]float64, len(rows))
	for i := range rows {
		row, ok := rows[i].([]interface{})
		fmt.Printf("%T\n", rows[i])
		fmt.Println(rows[i], ok)
		if !ok {
			return nil, false
		}
		a[i] = make([]float64, len(row))
		for j := range row {
			a[i][j], ok = row[j].(float64)
			fmt.Println(ok)
			if !ok {
				return nil, false
			}
		}
	}
	return a, true
}

type StretchGrayImage struct {
	Pix    [][]uint8
	Width  int
	Height int
}

func (img StretchGrayImage) At(x, y int) color.Color {
	pix := img.Pix[len(img.Pix)*y/img.Height][len(img.Pix[0])*x/img.Width]
	return color.RGBA{0, pix, pix, 255}
}

func (img StretchGrayImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

func (img StretchGrayImage) ColorModel() color.Model {
	return color.RGBAModel
}

func writeHtmlImg(w http.ResponseWriter, img image.Image) error {

	_, err := w.Write([]byte(`<img src="data:image/png;base64,`))
	if err != nil {
		return err
	}
	b64writer := base64.NewEncoder(base64.StdEncoding, w)
	err = png.Encode(b64writer, img)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(`" />`))
	if err != nil {
		return err
	}
	return nil
}

func minMax(a []float64) (float64, float64) {
	if len(a) == 0 {
		return 0, 0
	}
	min, max := a[0], a[0]
	for _, v := range a {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}
