package compose

import (
	"git.yo2.cz/drahoslav/penego/draw"
	"math"
)

func hitPlace(x, y float64, pos draw.Pos) bool {
	return math.Abs(pos.X-x) < 27 && math.Abs(pos.Y-y) < 27
}

func hitTransition(x, y float64, pos draw.Pos) bool {
	return math.Abs(pos.X-x) < 12 && math.Abs(pos.Y-y) < 38
}

func snap(x, y, n float64) draw.Pos {
	return draw.Pos{x - math.Mod(x, n), y - math.Mod(y, n)}
}
