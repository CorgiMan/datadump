# Datadump

Datadump prints images, plots, values and instances of structs directly to the browser. Datadump allows for quick prototyping of your application and a more advanced way to debug than the console.

## Example
[example output](example.html)

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
