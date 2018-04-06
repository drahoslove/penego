package compose

import (
	"math"

	"git.yo2.cz/drahoslav/penego/draw"
)

func hitPlace(x, y float64, pos draw.Pos) bool {
	r := draw.PLACE_RADIUS
	return math.Abs(pos.X-x) < r && math.Abs(pos.Y-y) < r
}

func hitTransition(x, y float64, pos draw.Pos) bool {
	rv, rh := draw.TRANSITION_WIDTH/2, draw.TRANSITION_HEIGHT/2
	return math.Abs(pos.X-x) < rv && math.Abs(pos.Y-y) < rh
}

func snap(x, y, n float64) draw.Pos {
	return draw.Pos{x - math.Mod(x, n), y - math.Mod(y, n)}
}

func center(ax, ay, bx, by float64) (x, y float64) {
	x = ax + (bx-ax)/2.0
	y = ay + (by-ay)/2.0
	return
}
