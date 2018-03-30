package compose

import (
	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)

type Composition struct {
	places      map[*net.Place]draw.Pos
	transitions map[*net.Transition]draw.Pos
	arcsIn      map[*net.Arc][2]draw.Pos
	arcsOut     map[*net.Arc][2]draw.Pos
}

func (comp Composition) DrawWith(drawer draw.Drawer) {
	for place, pos := range comp.places {
		drawer.DrawPlace(pos, place.Tokens, place.Description)
	}

	for tran, pos := range comp.transitions {
		drawer.DrawTransition(pos, tran.TimeFunc.String(), tran.Description)
	}

	for arc, pos := range comp.arcsIn {
		drawer.DrawInArc(pos[0], pos[1], arc.Weight)
	}

	for arc, pos := range comp.arcsOut {
		drawer.DrawOutArc(pos[0], pos[1], arc.Weight)
	}
}

func NewComposition() Composition {
	return Composition{
		make(map[*net.Place]draw.Pos),
		make(map[*net.Transition]draw.Pos),
		make(map[*net.Arc][2]draw.Pos),
		make(map[*net.Arc][2]draw.Pos),
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
		// arcs:
		for pi, p := range places {
			for _, arc := range tran.Origins {
				if arc.Place == p {
					composition.arcsIn[arc] = [2]draw.Pos{posOfPlace(pi), posOfTransition(ti)}
				}
			}
			for _, arc := range tran.Targets {
				if arc.Place == p {
					composition.arcsOut[arc] = [2]draw.Pos{posOfTransition(ti), posOfPlace(pi)}
				}
			}
		}
	}

	return composition
}
