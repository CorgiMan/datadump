package datadump

import (
	"fmt"
	"math/rand"
	"net/http"
)

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
