package main

import (
	"github.com/CorgiMan/datadump"
	"github.com/CorgiMan/jsonquery"
)

func main() {
	datadump.Open(":8080")
	defer datadump.Close()
	// // plot location markers on a map
	// geo := jsonquery.
	//  FromURL("https://ckannet-storage.commondatastorage.googleapis.com/2015-01-02T17:43:10.682Z/locations.json").
	//  Select(`{"latitude":"", "longitude":""}`).
	//  Flatten().Rename("latitude", "lat", "longitude", "lng")
	// for i := range geo["lat"] {
	//  geo["lat"][i], _ = strconv.ParseFloat(geo["lat"][i].(string), 64)
	// }
	// for i := range geo["lng"] {
	//  geo["lng"][i], _ = strconv.ParseFloat(geo["lng"][i].(string), 64)
	// }
	// datadump.Show(geo)
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
	datadump.Show(locs.Flatten())

	// extract
	plot := jsonquery.
		FromURL("http://www.asterank.com/api/skymorph/search?target=J99TS7A").
		Select(`{"pixel_loc_x":"", "pixel_loc_y":""}`).
		Flatten().Rename("pixel_loc_x", "x", "pixel_loc_y", "y")
	datadump.Show(plot)

	// // plot the json associated to the data
	jsondata := jsonquery.FromURL("http://www.asterank.com/api/skymorph/search?target=J99TS7A")
	datadump.Show(jsondata)
}
