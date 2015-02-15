package datadump

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"

	"github.com/CorgiMan/json2"
)

var C chan interface{}

func Open() chan interface{} {
	C = make(chan interface{})
	http.HandleFunc("/", root)
	http.HandleFunc("/ajax", ajaxdump(C))
	go func() { http.ListenAndServe(":8000", nil) }()
	OpenInBrowser("http://localhost:8000")
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
		_, err = fmt.Fprintf(w, "%s", x)
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
	var ys []interface{}
	if ys, ok = ValueArray(m, "y"); !ok {
		return
	}

	var xs []interface{}
	if xs, ok = ValueArray(m, "x"); !ok || len(ys) != len(xs) {
		for i := range ys {
			xs = append(xs, float64(i))
		}
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
			fmt.Fprintf(w, "%v,%v ", 300+8*xs[i].(float64), 175-8*ys[i].(float64))
		}
		fmt.Fprintf(w, `"style="fill:none;stroke:black;stroke-width:3" />`)
		return true, nil
	}
	if _, ok := m["closed"]; ok {
		//closed means polygon
		fmt.Fprintf(w, `<polygon points="`)
		for i := range ys {
			fmt.Fprintf(w, "%v,%v ", 300+8*xs[i].(float64), 175-8*ys[i].(float64))
		}
		fmt.Fprintf(w, `"style="fill:none;stroke:black;stroke-width:3" />`)
		return true, nil
	}
	if _, ok := m["lines"]; ok {
		//lines
		for i := 0; i+1 < len(ys); i += 2 {
			fmt.Fprintf(w, `<line x1="%v" y1="%v" x2="%v" y2="%v" 
				style="stroke:rgb(255,0,0);stroke-width:2" />`,
				300+8*xs[i].(float64), 175-8*ys[i].(float64),
				300+8*xs[i+1].(float64), 175-8*ys[i+1].(float64))
		}

		return true, nil
	}

	// write circles
	for i := range ys {
		fmt.Fprintf(w, `<circle cx="%v" cy="%v" r="%v" fill="red" />`,
			300+8*xs[i].(float64), 175-8*ys[i].(float64), r)
	}

	// color

	// widths, heights for rect
	return true, nil
}

func ValueArray(m map[string]interface{}, s string) (xs []interface{}, ok bool) {
	var a interface{}
	if a, ok = m[s]; !ok {
		return
	}

	if xs, ok = a.([]interface{}); !ok {
		return
	}

	//check if xs contains only values

	return xs, true
}

func writeGeoMap(w http.ResponseWriter, v interface{}) (bool, error) {
	return false, nil
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

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
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
