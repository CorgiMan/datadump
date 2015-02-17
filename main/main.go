package main

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"
	"os"
	"strconv"

	"github.com/CorgiMan/datadump"
	"github.com/CorgiMan/jsonquery"
)

func main() {
	datadump.Open(":8080")
	defer datadump.Close()

	// print a string
	datadump.C <- "Hello World!"

	// print a type (a file for example)
	f, _ := os.Open("main.go")
	datadump.C <- f

	// draw an image
	datadump.C <- Mandelbrot{Width: 300, Height: 300}

	// plot a sine function
	ys := []float64{}
	xs := []float64{}
	for x := -30.0; x < 30.0; x += 0.1 {
		xs = append(xs, x)
		ys = append(ys, math.Sin(x))
	}
	datadump.C <- map[string]interface{}{"connected": 0, "x": xs, "y": ys}

	// plot location markers on a map
	geo := jsonquery.
		FromURL("https://ckannet-storage.commondatastorage.googleapis.com/2015-01-02T17:43:10.682Z/locations.json").
		Select(`{"latitude":"", "longitude":""}`).
		Flatten().Rename("latitude", "lat", "longitude", "lng")
	for i := range geo["lat"] {
		geo["lat"][i], _ = strconv.ParseFloat(geo["lat"][i].(string), 64)
	}
	for i := range geo["lng"] {
		geo["lng"][i], _ = strconv.ParseFloat(geo["lng"][i].(string), 64)
	}
	datadump.C <- geo

	// plot location markers on a map
	locs := jsonquery.
		FromURL("http://www.amsterdamopendata.nl/files/ivv/parkeren/locaties.json").
		Select(`{"Locatie": ""}`)

	for i := range locs {
		m, ok := locs[i].(map[string]interface{})
		if !ok {
			continue
		}
		locs[i] = jsonquery.FromString(m["Locatie"].(string))[0]
	}
	datadump.C <- locs.Flatten()

	// plot coordinates extracted from a large json string
	plot := jsonquery.
		FromURL("http://www.asterank.com/api/skymorph/search?target=J99TS7A").
		Select(`{"pixel_loc_x":"", "pixel_loc_y":""}`).
		Flatten().Rename("pixel_loc_x", "x", "pixel_loc_y", "y")
	datadump.C <- plot

	// plot the json associated to the data
	jsondata := jsonquery.FromURL("http://www.asterank.com/api/skymorph/search?target=J99TS7A")
	datadump.C <- jsondata

}

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
