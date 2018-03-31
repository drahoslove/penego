package compose

import (
	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)

type Composable interface {
}

type Composition struct {
	places      map[*net.Place]draw.Pos
	transitions map[*net.Transition]draw.Pos
}

func NewComposition() Composition {
	return Composition{
		make(map[*net.Place]draw.Pos),
		make(map[*net.Transition]draw.Pos),
	}
}

func (comp Composition) HitTest(x, y float64) Composable {
	for place, pos := range comp.places {
		if hitPlace(x, y, pos) {
			return place
		}
	}
	for transition, pos := range comp.transitions {
		if hitTransition(x, y, pos) {
			return transition
		}
	}
	return nil
}

func (comp Composition) Move(node Composable, x, y float64) {
	pos := draw.Pos{x, y}
	switch v := node.(type) {
	case *net.Transition:
		comp.transitions[v] = pos
	case *net.Place:
		comp.places[v] = pos
	}
}

func (comp Composition) GhostMove(node Composable, x, y float64) {
	println("GhostMove Not implemented")
}

func (comp Composition) DrawWith(drawer draw.Drawer) {
	for place, pos := range comp.places {
		drawer.DrawPlace(pos, place.Tokens, place.Description)
	}

	for tran, pos := range comp.transitions {
		drawer.DrawTransition(pos, tran.TimeFunc.String(), tran.Description)
		for _, arc := range tran.Origins {
			from := comp.places[arc.Place]
			drawer.DrawInArc(from, pos, arc.Weight)
		}
		for _, arc := range tran.Targets {
			to := comp.places[arc.Place]
			drawer.DrawOutArc(pos, to, arc.Weight)
		}
	}
}

// basic "dumb" way to draw a net
func GetSimple(network net.Net) Composition {
	places := network.Places()
	transitions := network.Transitions()

	const BASE = 90.0

	posOfPlace := func(i int) draw.Pos {
		pos := draw.Pos{
			X: float64(i)*BASE - (float64(len(places))/2-0.5)*BASE,
			Y: 0,
		}
		if len(transitions) <= 1 {
			pos.Y += BASE
		}
		return pos
	}

	posOfTransition := func(i int) draw.Pos {
		pos := draw.Pos{
			X: float64(i)*BASE - (float64(len(transitions))/2)*BASE + BASE/2,
			Y: 4*BASE*float64(i%2) - 2*BASE,
		}
		if len(transitions) <= 1 {
			pos.Y += BASE
		}
		return pos
	}

	composition := NewComposition()

	// compute positions of places
	for i, place := range places {
		composition.places[place] = posOfPlace(i)
	}
	// compute positions of transitions
	for ti, tran := range transitions {
		composition.transitions[tran] = posOfTransition(ti)
	}

	return composition
}
