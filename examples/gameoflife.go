package main

import (
	"time"

	"github.com/CorgiMan/datadump"
	"github.com/CorgiMan/jsonquery"
)

type Cell struct {
	x, y int
}

var cells map[Cell]bool

func (c Cell) neighs() []Cell {
	return []Cell{
		{c.x + 1, c.y}, {c.x + 1, c.y + 1}, {c.x, c.y + 1}, {c.x - 1, c.y + 1},
		{c.x - 1, c.y}, {c.x - 1, c.y - 1}, {c.x, c.y - 1}, {c.x + 1, c.y - 1}}
}

func (c Cell) active_neighs() []Cell {
	res := c.neighs()
	for i := 0; i < len(res); i++ {
		if !cells[res[i]] {
			res = append(res[:i], res[i+1:]...)
			i--
		}
	}
	return res
}

func step() {
	next := make(map[Cell]bool)
	for c := range cells {
		for _, cc := range c.neighs() {
			n := len(cc.active_neighs())
			if n == 3 || (n == 2 && cells[cc]) {
				next[cc] = true
			}
		}
	}
	cells = next
}

func main() {
	datadump.Open(":8080")
	defer datadump.Close()
	cells = make(map[Cell]bool)
	for x := -3; x <= 6; x++ {
		cells[Cell{x, 0}] = true
	}
	for {
		keys := []interface{}{}
		for c := range cells {
			keys = append(keys, c)
		}
		m := jsonquery.From(keys).Select(`{"x":"","y":""}`).Flatten()
		m["x"], m["y"] = append(m["x"], -10), append(m["y"], 10)
		m["x"], m["y"] = append(m["x"], 10), append(m["y"], -10)
		datadump.Show(m)
		step()
		time.Sleep(time.Second * 2 / 10)
	}
}
