# Datadump

Datadump prints images, plots, values and instances of structs directly to the browser. Datadump allows for quick prototyping of your application and a more advanced way to debug than the console.

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

## Example
This [example output](http://rawgit.com/CorgiMan/datadump/master/example.html) is a result of running the code in the usage section.

## Usage
The api consists of 3 elements: a `chan interface{}` called `datadump.C` in which we sent the data to be shown in the browser, together with 2 functions `datadump.Open(port string)` and `datadump.Close()`.

First we need to open a port to sent your data to. 
```
datadump.Open(":8080")
defer datadump.Close()
```

Next we can sent stuff to the `datadump.C` channel to sent stuff to the browser.

```
// Print a string
datadump.Show("Hello World!")
```
 
``` 
// Pretty print the contents of a file struct
f, _ := os.Open("main.go")
datadump.Show(f)
```

```
// Show an Image (the mandelbrot image is defined in main/main.go)
datadump.Show(Mandelbrot{Width: 300, Height: 300})
```

```
// Plot a sin function
// xs and ys are of type []float64 and are defined in main/main.go
datadump.Show(map[string]interface{}{"connected": 0, "x": xs, "y": ys})
```

```
// Plot datapoints from data from the web. We open a json file 
// from the web and then transform it using jsonquery
datadump.Show(jsonquery.FromURL("http://www.asterank.com/api/skymorph/search?target=J99TS7A").)
                        Select(`{"pixel_loc_x":"", "pixel_loc_y":""}`).
                        Flatten().
                        Rename("pixel_loc_x", "x", "pixel_loc_y", "y")
```

## Installation
You need the `json2` package. The `jsonquery` is optional but synergizes well with this application, as shown in the example.

```
go get github.com/CorgiMan/datadump
go get github.com/CorgiMan/jsonquery
go get github.com/CorgiMan/json2
```

`json2` is exactly the same as encoding/json but it also marshals unexported fields. This is used to pretty print any type as a json string.

`jsonquery` is used in the last example. It selects all the objects that match a given json string.

run the example with `cd main` and then `go run main.go`
