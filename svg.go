package datadump

import (
	"fmt"
	"net/http"
)

func writeSvg(w http.ResponseWriter, v interface{}) (ok bool, err error) {
	var m map[string]interface{}
	if m, ok = v.(map[string]interface{}); !ok {
		return
	}
	var ys []float64
	if ys, ok = ValueArray(m, "y"); !ok {
		return
	}

	var xs []float64
	if xs, ok = ValueArray(m, "x"); !ok || len(ys) != len(xs) {
		for i := range ys {
			xs = append(xs, float64(i))
		}
	}

	minx, maxx := minMax(xs)
	miny, maxy := minMax(ys)
	minx -= 0.05 * (maxx - minx)
	maxx += 0.05 * (maxx - minx)
	miny -= 0.05 * (maxy - miny)
	maxy += 0.05 * (maxy - miny)
	for i := range xs {
		xs[i] = 600 * (xs[i] - minx) / (maxx - minx)
		ys[i] = 350 - 350*(ys[i]-miny)/(maxy-miny)
	}

	var r interface{}
	if r, ok = m["r"]; !ok {

		r = 3.0
	}

	w.Write([]byte(`<svg width="600" height="350">`))
	defer w.Write([]byte(`</svg>`))

	if _, ok := m["connected"]; ok {

		fmt.Fprintf(w, `<polyline points="`)
		for i := range ys {
			fmt.Fprintf(w, "%v,%v ", xs[i], ys[i])
		}
		fmt.Fprintf(w, `"style="fill:none;stroke:black;stroke-width:3" />`)
		return true, nil
	}
	if _, ok := m["closed"]; ok {

		fmt.Fprintf(w, `<polygon points="`)
		for i := range ys {
			fmt.Fprintf(w, "%v,%v ", xs[i], ys[i])
		}
		fmt.Fprintf(w, `"style="fill:none;stroke:black;stroke-width:3" />`)
		return true, nil
	}
	if _, ok := m["lines"]; ok {

		for i := 0; i+1 < len(ys); i += 2 {
			fmt.Fprintf(w, `<line x1="%v" y1="%v" x2="%v" y2="%v" 
                style="stroke:rgb(255,0,0);stroke-width:2" />`,
				xs[i], ys[i], xs[i+1], ys[i+1])
		}

		return true, nil
	}

	for i := range ys {
		fmt.Fprintf(w, `<circle cx="%v" cy="%v" r="%v" fill="red" />`,
			xs[i], ys[i], r)
	}

	return true, nil
}

func ValueArray(m map[string]interface{}, s string) ([]float64, bool) {
	var a interface{}
	var ok bool
	var xs []interface{}
	if a, ok = m[s]; !ok {
		return nil, false
	}

	if xs, ok = a.([]interface{}); !ok {
		return nil, false
	}

	//check if xs contains only values

	ys := []float64{}
	for i := range xs {
		// if s, ok := xs[i].(string); ok {
		// 	if x, okorerr := strconv.ParseFloat(s, 64); okorerr!===nil {
		// 		ys = append(ys, x)
		// 	}
		// } else
		if x, ok := xs[i].(float64); ok {
			ys = append(ys, x)
		}
	}
	return ys, true
}
