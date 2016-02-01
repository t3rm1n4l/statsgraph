package main

import (
	"github.com/dustin/go-humanize"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"image/color"
)

var colors = []color.RGBA{
	color.RGBA{R: 53, G: 142, B: 242, A: 255},
	color.RGBA{R: 255, G: 215, B: 0, A: 255},
	color.RGBA{R: 220, G: 20, B: 60, A: 255},
	color.RGBA{R: 255, G: 165, B: 0, A: 255},
	color.RGBA{R: 93, G: 65, B: 75, A: 255},
	color.RGBA{R: 0, G: 128, B: 0, A: 255},
	color.RGBA{R: 138, G: 43, B: 226, A: 255},
	color.RGBA{R: 255, G: 0, B: 0, A: 255},
}

type ticks struct {
}

func (tt *ticks) Ticks(min, max float64) (t []plot.Tick) {
	x := (max - min) / 10
	for i := 0; i < 10; i++ {
		t = append(t, plot.Tick{Value: min + x*float64(i), Label: humanize.Commaf(min + x*float64(i))})
	}
	return
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

	p.p.X.Min = 0
	p.p.Y.Min = 0
	p.p.Y.Padding = 0
	p.p.X.Padding = 0
	p.p.Y.Tick.Marker = &ticks{}

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
		l.LineStyle.Color = colors[i%len(colors)]
		p.p.Legend.Add(ln.name, l)
		p.p.Add(l)
	}

	return p.p.Save(10*vg.Inch, 10*vg.Inch, f)
}
