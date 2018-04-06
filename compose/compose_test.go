package compose

import (
	"testing"

	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)

func getComp() Composition {
	comp := New()
	comp.transitions[&net.Transition{}] = draw.Pos{10, 25}
	comp.transitions[&net.Transition{}] = draw.Pos{30, 45}
	comp.places[&net.Place{}] = draw.Pos{40, 15}
	comp.places[&net.Place{}] = draw.Pos{30, 25}

	return comp
}

func TestCompositionFindCenter(test *testing.T) {
	comp := getComp()

	centerX, centerY := comp.FindCenter()
	if centerX != 25 {
		test.Errorf("value centerX should be %f not %f", 25., centerX)
	}
	if centerY != 30 {
		test.Errorf("Value centerY should be %f not %f", 30., centerY)
	}
}

func TestCompositionCenterTo(test *testing.T) {
	comp := getComp()

	{
		// center to zero
		x, y := 0.0, 0.0
		comp.CenterTo(x, y)
		centerX, centerY := comp.FindCenter()
		if centerX != x || centerY != y {
			test.Errorf("center should be %v;%v not %v;%v", x, y, centerX, centerY)
		}
	}
	{
		// center to
		x, y := 15., 45.
		comp.CenterTo(x, y)
		centerX, centerY := comp.FindCenter()
		if centerX != x || centerY != y {
			test.Errorf("center should be %v;%v not %v;%v", x, y, centerX, centerY)
		}
	}
	{
		// center to
		x, y := 17., 49.
		pos := snap(x, y, 15)
		x = pos.X
		y = pos.Y
		comp.CenterTo(x, y)
		centerX, centerY := comp.FindCenter()
		if centerX != x || centerY != y {
			test.Errorf("center should be %v;%v not %v;%v", x, y, centerX, centerY)
		}
	}
}
