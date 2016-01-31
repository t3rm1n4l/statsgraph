package main

import (
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"image/color"
)

var colors = []color.RGBA{
	color.RGBA{R: 53, G: 142, B: 242, A: 255},
	color.RGBA{R: 40, G: 48, B: 68, A: 255},
	color.RGBA{R: 120, G: 161, B: 187, A: 255},
	color.RGBA{R: 139, G: 120, B: 109, A: 255},
	color.RGBA{R: 93, G: 65, B: 75, A: 255},
	color.RGBA{R: 80, G: 6, B: 36, A: 255},
}

type Plot struct {
	p     *plot.Plot
	lines []*Line
}

type Line struct {
	name string
	pts  plotter.XYs
}

func NewPlot(title string) *Plot {
	p0, _ := plot.New()
	p := &Plot{
		p: p0,
	}

	p.p.Title.Text = title
	p.p.X.Label.Text = "time"
	p.p.Add(plotter.NewGrid())

	return p
}

func (p *Plot) NewLine(s string) *Line {
	line := &Line{name: s}
	p.lines = append(p.lines, line)
	return line
}

func (line *Line) AddPoint(x, y float64) {
	var pt struct{ X, Y float64 }
	pt.X = x
	pt.Y = y
	line.pts = append(line.pts, pt)
}

func (p *Plot) Write(f string) error {
	for i, ln := range p.lines {
		l, _ := plotter.NewLine(ln.pts)
		l.LineStyle.Width = vg.Points(2)
		l.LineStyle.Color = colors[i]
		p.p.Legend.Add(ln.name, l)
		p.p.Add(l)
	}

	return p.p.Save(10*vg.Inch, 10*vg.Inch, f)
}
