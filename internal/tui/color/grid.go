package color

import (
	"math/rand"
	"sync/atomic"

	"github.com/lucasb-eyer/go-colorful"
)

var PbColors = newPbColor()

type PbColor struct {
	idx    *atomic.Int32
	colors [][]string
}

const (
	defaultColorX = 14
	defaultColorY = 8
)

func newPbColor() PbColor {
	i := &atomic.Int32{}
	// no need to use secure random generator
	i.Store(rand.Int31()) //nolint:gosec

	return PbColor{
		idx:    i,
		colors: colorGrid(defaultColorX, defaultColorY),
	}
}

func (p *PbColor) Colors() []string {
	p.idx.Add(1)
	i := p.idx.Load() % 8 //nolint:gomnd
	return p.colors[i]
}

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, 2) //nolint:gomnd
		grid[x][0] = y0.BlendLuv(x1[x], float64(0)/float64(xSteps)).Hex()
		grid[x][1] = y0.BlendLuv(x1[x], float64(xSteps-1)/float64(xSteps)).Hex()
	}

	return grid
}
