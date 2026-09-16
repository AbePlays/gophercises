package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
)

func main() {
	canvas := svg.New(os.Stdout)

	data := []struct {
		Month string
		Usage int
	}{
		{"Jan", 171},
		{"Feb", 180},
		{"Mar", 134},
		{"Apr", 200},
		{"May", 150},
		{"Jun", 160},
		{"Jul", 87},
		{"Aug", 169},
		{"Sep", 103},
		{"Oct", 140},
		{"Nov", 150},
		{"Dec", 160},
	}

	w, h := len(data)*60+10, 300
	maxHeight := 0
	threshold := 120
	for _, item := range data {
		if item.Usage > maxHeight {
			maxHeight = item.Usage
		}
	}

	canvas.Start(w, h)
	for i, v := range data {
		percent := v.Usage * (h - 50) / maxHeight
		canvas.Rect(i*60+10, (h-50)-percent, 50, percent, "fill:rgb(77,200,232)")
		canvas.Text(i*60+35, h-25, v.Month, "font-size:16px;text-anchor:middle;fill:rgb(150,150,150)")
	}

	threshPercent := threshold * (h - 50) / maxHeight
	canvas.Line(0, h-threshPercent, w, h-threshPercent, "stroke:rgb(250,100,100);stroke-width:2px;")
	canvas.Line(0, h-50, w, h-50, "stroke:rgb(150,150,150);stroke-width:2px;")
	canvas.Rect(0, 0, w, h-threshPercent, "fill:rgb(255, 100, 100);opacity:0.3")

	canvas.End()
}
