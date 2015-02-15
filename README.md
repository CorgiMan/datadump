# Datadump

Datadump prints images, plots, values and instances of structs directly to the browser. Datadump allows for quick prototyping of your application and a more advanced way to debug than the console.

## Example
[example output](http://rawgit.com/CorgiMan/datadump/master/example.html)

## Usage
The example above is a result of the following code

You need to open a port to sent your data to. 
```
datadump.Open(":8080")
defer datadump.Close()
```

Sent your data to the `datadump.C` channel
- print a string
`datadump.C <- "Hello World!"`
 
- print the contents of a file
``` 
f, _ := os.Open("main.go")
datadump.C <- f
```

- show an Image (the mandelbrot image is defined in main/main.go)
```
datadump.C <- Mandelbrot{Width: 300, Height: 300}
```

- Plot a sin function (xs and ys are of type []float64 and are defined in main/main.go)
```
datadump.C <- map[string]interface{}{"connected": 0, "x": xs, "y": ys}
```

- Plot some datapoints found in a json file from the web. The json is transformed to a graphable form with jsonquery.
```
datadump.C <- FromURL("http://www.asterank.com/api/skymorph/search?target=J99TS7A").
              Select(`{"pixel_loc_x":"", "pixel_loc_y":""}`).
              Flatten().
              Rename("pixel_loc_x", "x", "pixel_loc_y", "y")
```

## Features
- draw images
- graph circles, rectangles, points, polygons
- draw an instance of a struct as a json string

## Future features
- plot geo coordinates on a map
- audio support
- plots with axis
- 3d plots
- graphs and tree visualization
- support for cyclic types

## Installation
You need the `json2` package. The `jsonquery` is optional but synergizes well with this application, as shown in the example `main/main.go`.

```
go get github.com/CorgiMan/datadump
go get github.com/CorgiMan/jsonquery
go get github.com/CorgiMan/json2
```

`jsonquery` is used in the last example. It selects all the objects that match a given json string.

`json2` is exactly the same as encoding/json but it also marshals unexported fields. This is used to pretty print any type as a json string.

run example with `cd main` and then `go run main.go`
