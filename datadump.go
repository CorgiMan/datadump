package datadump

import (
	"fmt"
	"html/template"
	"image"
	"log"
	"net/http"

	"github.com/CorgiMan/json2"
)

var c chan interface{}

func Show(v interface{}) {
	c <- v
}

func Open(port string) chan interface{} {
	c = make(chan interface{})
	http.HandleFunc("/", root)
	http.HandleFunc("/ajax", ajaxdump(c))
	go func() { http.ListenAndServe(port, nil) }()
	OpenInBrowser("http://localhost" + port)
	return c
}

func Close() {
	c <- ""
}

type Handler func(http.ResponseWriter, *http.Request)

func root(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, nil)
}

func ajaxdump(c chan interface{}) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if ok, err = write2DImage(w, o); ok {
			break
		}

		_, err = fmt.Fprintf(w, "<pre>%s</pre>", bts)
	}
	if err != nil {
		log.Print(err)
	}
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<script src="http://maps.googleapis.com/maps/api/js"> </script>
<script>
var scripts = dpcument.body.getElementsByTagName("script");
for( var i=0; i<scripts.length; i++ ) {
    eval(scripts[i].innerText);
}
</script>

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
