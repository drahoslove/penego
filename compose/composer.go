// This package implements Composition
// which is a stuct defining 2-dimensional placeme of petri net elements
package compose

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/net"
)

type path struct {
	from Composable
	to   Composable
}

type Composable interface {
}

type Composition struct {
	places      map[*net.Place]draw.Pos
	transitions map[*net.Transition]draw.Pos
	pathes      map[*path][]draw.Pos
	ghosts      map[Composable]draw.Pos
}

func New() Composition {
	return Composition{
		make(map[*net.Place]draw.Pos),
		make(map[*net.Transition]draw.Pos),
		make(map[*path][]draw.Pos),
		make(map[Composable]draw.Pos),
	}
}

func (comp Composition) PathPositions(from Composable, to Composable) []draw.Pos {
	getPos := func(node Composable) draw.Pos {
		switch node := node.(type) {
		case *net.Place:
			return comp.places[node]
		case *net.Transition:
			return comp.transitions[node]
		default:
			return draw.Pos{}
		}
	}
	fromPos, toPos := getPos(from), getPos(to)

	positions := []draw.Pos{}

	// insert intermediate path positions
	for path, poss := range comp.pathes {
		if path.from == from && path.to == to {
			for _, pos := range poss {
				positions = append(positions, pos)
			}
		}
		if path.from == to && path.to == from {
			for _, pos := range poss {
				positions = append([]draw.Pos{pos}, positions...)
			}
		}
	}

	positions = append([]draw.Pos{fromPos}, positions...)
	positions = append(positions, toPos)
	return positions
}

func (comp Composition) String() string {
	str := ""
	for place, pos := range comp.places {
		str += fmt.Sprintf("%s %v;%v\n", place.Id, pos.X, pos.Y)
	}
	str += fmt.Sprintf("----\n")
	for transition, pos := range comp.transitions {
		str += fmt.Sprintf("%s %v;%v\n", transition.Id, pos.X, pos.Y)
	}
	return str
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
	pos := snap(x, y, 15)
	switch node := node.(type) {
	case *net.Transition:
		comp.transitions[node] = pos
	case *net.Place:
		comp.places[node] = pos
	}
	delete(comp.ghosts, node)
}

func (comp Composition) FindCenter() (float64, float64) {
	left, right := math.Inf(+1), math.Inf(-1)
	top, bottom := math.Inf(+1), math.Inf(-1)
	/*
	                -
	               top
	              ____
	   - left -> |    | <- right +
	             |____|
	             bottom
	               +
	*/
	enhanceEdges := func(pos draw.Pos) {
		x, y := pos.X, pos.Y
		if x < left {
			left = x
		}
		if x > right {
			right = x
		}
		if y < top {
			top = y
		}
		if y > bottom {
			bottom = y
		}
	}

	for _, pos := range comp.places {
		enhanceEdges(pos)
	}

	for _, pos := range comp.transitions {
		enhanceEdges(pos)
	}

	return center(left, top, right, bottom)
}

// move whole composition so its center is at x, y
// note that is uses Movew which snaps to multiples of 15, so it might not end up exactly on those positions
func (comp Composition) CenterTo(x, y float64) {
	centerX, centerY := comp.FindCenter()

	deltaX := x - centerX
	deltaY := y - centerY

	for place, pos := range comp.places {
		comp.Move(place, pos.X+deltaX, pos.Y+deltaY)
	}

	for tran, pos := range comp.transitions {
		comp.Move(tran, pos.X+deltaX, pos.Y+deltaY)
	}

	for _, poss := range comp.pathes {
		for i, _ := range poss {
			poss[i].X += deltaX
			poss[i].Y += deltaY
		}
	}
}

func (comp Composition) GhostMove(node Composable, x, y float64) {
	comp.ghosts[node] = snap(x, y, 15)
}

func (comp Composition) DrawWith(drawer draw.Drawer) {
	setStyle := func(node Composable) func(Composable) {
		if _, isGhosted := comp.ghosts[node]; isGhosted {
			drawer.SetStyle(draw.FadedStyle)
			return func(Composable) {
				drawer.SetStyle(draw.FadedStyle)
			}
		} else {
			drawer.SetStyle(draw.DefaultStyle)
			return func(node Composable) {
				if _, isGhosted := comp.ghosts[node]; isGhosted {
					drawer.SetStyle(draw.FadedStyle)
				} else {
					drawer.SetStyle(draw.DefaultStyle)
				}
			}
		}
	}

	// draw arcs
	for tran, _ := range comp.transitions {
		orSetStyle := setStyle(tran)
		for _, arc := range tran.Origins {
			if arc.Place.Hidden() {
				continue
			}
			orSetStyle(arc.Place)
			if arc.Type == net.InhibitorArc {
				drawer.DrawInhibitorArc(comp.PathPositions(arc.Place, tran))
			} else {
				drawer.DrawInArc(comp.PathPositions(arc.Place, tran), arc.Weight)
			}
		}
		for _, arc := range tran.Targets {
			if arc.Place.Hidden() {
				continue
			}
			orSetStyle(arc.Place)
			drawer.DrawOutArc(comp.PathPositions(tran, arc.Place), arc.Weight)
		}
	}

	// draw all places
	for place, pos := range comp.places {
		setStyle(place)
		drawer.DrawPlace(pos, place.Tokens, place.Description)
	}

	// draw all transtitions
	for tran, pos := range comp.transitions {
		setStyle(tran)
		drawer.DrawTransition(pos, tran.TimeFunc.String(), tran.Description)
	}

	// draw moving items last
	drawer.SetStyle(draw.HighlightedStyle)
	for node, pos := range comp.ghosts {
		switch node := node.(type) {
		case *net.Place:
			place := node
			drawer.DrawPlace(pos, place.Tokens, place.Description)
			for tran, _ := range comp.transitions {
				for _, arc := range tran.Origins {
					if arc.Place == place {
						if arc.Type == net.InhibitorArc {
							drawer.DrawInhibitorArc(comp.PathPositions(place, tran))
						} else {
							drawer.DrawInArc(comp.PathPositions(place, tran), arc.Weight)
						}
					}
				}
				for _, arc := range tran.Targets {
					if arc.Place == place {
						drawer.DrawOutArc(comp.PathPositions(tran, place), arc.Weight)
					}
				}
			}
		case *net.Transition:
			tran := node
			drawer.DrawTransition(pos, tran.TimeFunc.String(), tran.Description)
			for _, arc := range tran.Origins {
				if arc.Place.Hidden() {
					continue
				}
				if arc.Type == net.InhibitorArc {
					drawer.DrawInhibitorArc(comp.PathPositions(arc.Place, tran))
				} else {
					drawer.DrawInArc(comp.PathPositions(arc.Place, tran), arc.Weight)
				}
			}
			for _, arc := range tran.Targets {
				if arc.Place.Hidden() {
					continue
				}
				drawer.DrawOutArc(comp.PathPositions(node, arc.Place), arc.Weight)
			}
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

	composition := New()

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

func Parse(str string, network net.Net) Composition {
	composition := New()
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		id := parts[0]
		if len(parts) > 1 {
			poss := strings.Split(parts[1], ";")
			x, _ := strconv.ParseFloat(poss[0], 64)
			y, _ := strconv.ParseFloat(poss[1], 64)
			for _, tran := range network.Transitions() {
				if tran.Id == id {
					composition.transitions[tran] = draw.Pos{x, y}
				}
			}
			for _, place := range network.Places() {
				if place.Id == id {
					composition.places[place] = draw.Pos{x, y}
				}
			}
		}
	}
	return composition
}
