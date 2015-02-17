package datadump

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"

	"github.com/CorgiMan/json2"
)

var C chan interface{}

func Open(port string) chan interface{} {
	C = make(chan interface{})
	http.HandleFunc("/", root)
	http.HandleFunc("/ajax", ajaxdump(C))
	go func() { http.ListenAndServe(port, nil) }()
	OpenInBrowser("http://localhost" + port)
	return C
}

func Close() {
	C <- ""
}

type Handler func(http.ResponseWriter, *http.Request)

func root(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, nil)
}

func ajaxdump(c chan interface{}) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// v := <-c
		writeInterface(w, <-c)
	}
}

func writeInterface(w http.ResponseWriter, v interface{}) {
	var err error
	switch x := v.(type) {
	case image.Image:
		err = writeHtmlImg(w, x)
	case string:
		_, err = fmt.Fprintf(w, "<pre>%s</pre>", x)
	default:
		// try to marshal and unmarshal and then find if interface contains
		// svg or geo elements
		var bts []byte
		bts, err = json2.MarshalIndent(x, "", "    ")
		if err != nil {
			_, err = fmt.Fprintf(w, "type %T: %v", x, x)
			break
		}
		var o interface{}
		json2.Unmarshal(bts, &o)

		var ok bool
		if ok, err = writeSvg(w, o); ok {
			break
		}
		if ok, err = writeGeoMap(w, o); ok {
			break
		}

		_, err = fmt.Fprintf(w, "<pre>%s</pre>", bts)
	}
	if err != nil {
		log.Print(err)
	}
}

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
		// also check is m["r"] is a value
		r = 3.0
	}

	w.Write([]byte(`<svg width="600" height="350">`))
	defer w.Write([]byte(`</svg>`))

	if _, ok := m["connected"]; ok {
		//connected means polyline
		fmt.Fprintf(w, `<polyline points="`)
		for i := range ys {
			fmt.Fprintf(w, "%v,%v ", xs[i], ys[i])
		}
		fmt.Fprintf(w, `"style="fill:none;stroke:black;stroke-width:3" />`)
		return true, nil
	}
	if _, ok := m["closed"]; ok {
		//closed means polygon
		fmt.Fprintf(w, `<polygon points="`)
		for i := range ys {
			fmt.Fprintf(w, "%v,%v ", xs[i], ys[i])
		}
		fmt.Fprintf(w, `"style="fill:none;stroke:black;stroke-width:3" />`)
		return true, nil
	}
	if _, ok := m["lines"]; ok {
		//lines
		for i := 0; i+1 < len(ys); i += 2 {
			fmt.Fprintf(w, `<line x1="%v" y1="%v" x2="%v" y2="%v" 
				style="stroke:rgb(255,0,0);stroke-width:2" />`,
				xs[i], ys[i], xs[i+1], ys[i+1])
		}

		return true, nil
	}

	// write circles
	for i := range ys {
		fmt.Fprintf(w, `<circle cx="%v" cy="%v" r="%v" fill="red" />`,
			xs[i], ys[i], r)
	}

	// color

	// widths, heights for rect
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
		ys = append(ys, xs[i].(float64))
	}
	return ys, true
}

func writeGeoMap(w http.ResponseWriter, v interface{}) (ok bool, err error) {
	var m map[string]interface{}
	if m, ok = v.(map[string]interface{}); !ok {
		return false, nil
	}

	coords, hascoords := m["coordinates"]
	lats, haslat := ValueArray(m, "lat")
	lngs, haslng := ValueArray(m, "lng")
	if !hascoords && !(haslat && haslng) {
		return false, nil
	}
	if hascoords {
		for _, c := range coords.([]interface{}) {
			lngs = append(lngs, c.([]interface{})[0].(float64))
			lats = append(lats, c.([]interface{})[1].(float64))
		}
	}

	x := rand.Intn(100000)

	fmt.Fprintf(w,
		`
		<script>
		  var mapProp = {
		    center: new google.maps.LatLng(51.508742,-0.120850),
		    zoom:9,
		    mapTypeId: google.maps.MapTypeId.ROADMAP
		  };
		  var map = new google.maps.Map(document.getElementById("%d"),mapProp);
		  var bounds = new google.maps.LatLngBounds();
        `, x)
	for i := range lats {
		fmt.Fprintf(w,

			`
			  loc = new google.maps.LatLng(%f, %f)
			  bounds.extend(loc);
			  var marker = new google.maps.Marker({
	      		position: loc,
	      		map: map,
	      		title: 'Hello World!'
	  		  });

			`, lats[i], lngs[i])
	}

	fmt.Fprintf(w, `

		    map.fitBounds(bounds);
		    map.panToBounds(bounds); 
		</script>
		`)

	fmt.Fprintf(w,
		`
		<div id="%d" style="width:600px;height:350px;"></div>
		`, x)

	return true, nil
}

func writeHtmlImg(w http.ResponseWriter, img image.Image) error {
	// encoded image gets written to b64writer, which writes to httpwriter
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

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<script src="http://maps.googleapis.com/maps/api/js"> </script>

<script>
function loadXMLDoc() {
    var xmlhttp;
    if (window.XMLHttpRequest) {
      // code for IE7+, Firefox, Chrome, Opera, Safari
      xmlhttp=new XMLHttpRequest();
    } else {
      // code for IE6, IE5
      xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
    }

    xmlhttp.onreadystatechange=function() {
      if (xmlhttp.readyState==4 && xmlhttp.status==200) {
        var node = document.createElement("div");
        node.innerHTML = xmlhttp.responseText;
        var outer = document.getElementById("myDiv")
        outer.appendChild(node);

var scripts = node.getElementsByTagName("script");
for( var i=0; i<scripts.length; i++ ) {
    eval(scripts[i].innerText);
}

        var node2 = document.createElement("div")
        node2.style.border = "1px solid black" 
        outer.appendChild(node2)
        xmlhttp.open("GET","/ajax",true);
        xmlhttp.send();
      }
    }

    xmlhttp.open("GET","/ajax",true);
    xmlhttp.send();
}
loadXMLDoc()

</script>
</head>

<body>
<h2>DataDump</h2>
<div id="myDiv" />
</body>
</html>
`))
