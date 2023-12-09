package color

import "github.com/lucasb-eyer/go-colorful"

var PbColors = colorGrid(14, 8)

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
		grid[x] = make([]string, 2)
		grid[x][0] = y0.BlendLuv(x1[x], float64(0)/float64(xSteps)).Hex()
		grid[x][1] = y0.BlendLuv(x1[x], float64(xSteps-1)/float64(xSteps)).Hex()
	}

	return grid
}
