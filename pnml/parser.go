// Package pnml implements parser of pnml format
// Currently supported dialects are those from PIPE 5 and CPN tools 4
// Not all features might be supported
package pnml

import (
	"encoding/xml"
	"io"
	"strconv"
	"strings"

	"git.yo2.cz/drahoslav/penego/compose"
	"git.yo2.cz/drahoslav/penego/net"
)

// structures defining pnml format

type Pnml struct {
	Net Net `xml:"net"`
}

type Net struct {
	Pages       []Page       `xml:"page"`
	Places      []Place      `xml:"place"`
	Transitions []Transition `xml:"transition"`
	Arcs        []Arc        `xml:"arc"`
}

type Page struct {
	Places      []Place      `xml:"place"`
	Transitions []Transition `xml:"transition"`
	Arcs        []Arc        `xml:"arc"`
}

type Place struct {
	Id       string   `xml:"id,attr"`
	Name     Val      `xml:"name"`
	Marking  Val      `xml:"initialMarking"`
	Position Position `xml:"graphics>position"`
}

type Position struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}

type Transition struct {
	Id       string   `xml:"id,attr"`
	Name     Val      `xml:"name"`
	Priority int      `xml:"priority>value"`
	Position Position `xml:"graphics>position"`
}

type Arc struct {
	Id      string `xml:"id,attr"`
	Source  string `xml:"source,attr"`
	Target  string `xml:"target,attr"`
	Type    Val    `xml:"type"`
	ArcType Val    `xml:"arctype"`
	Weight  Val    `xml:"inscription"`
}

// NOTE
//   PIPE uses <value>
//   CPN uses <text>
type Val struct {
	Text      string `xml:"text"`
	Value     string `xml:"value"`
	ValueAttr string `xml:"value,attr"`
}

func (v Val) String() string {
	if v.Text != "" {
		return v.Text
	}
	if v.Value != "" {
		return v.Value
	}
	return v.ValueAttr
}

func (v Val) Int(def int) int {
	if v.Text != "" {
		val, err := strconv.Atoi(v.Value)
		if err == nil {
			return val
		}
	}
	if v.Value != "" {
		val, err := strconv.Atoi(v.Value)
		if err == nil {
			return val
		}
		// PIPE has some values prefixed with Default,
		parts := strings.SplitAfter(v.Value, "Default,")
		if len(parts) == 2 {
			val, err = strconv.Atoi(parts[1])
			if err == nil {
				return val
			}
		}
	}
	return def
}
func (pnml *Pnml) buildNetCompo() (net.Net, compose.Composition) {
	composition := compose.New()

	places := net.Places{}
	transitions := net.Transitions{}

	buildPlaces := func(pnmlPlaces []Place) {
		for _, p := range pnmlPlaces {
			place := &net.Place{
				Tokens:      p.Marking.Int(0),
				Id:          p.Id,
				Description: p.Name.String(),
			}

			composition.Move(place, p.Position.X, p.Position.Y)
			places.Push(place)
		}
	}
	buildTransitions := func(pnmlTtransitions []Transition, pnmlArcs []Arc) {
		for _, t := range pnmlTtransitions {
			origins := net.Arcs{}
			targets := net.Arcs{}

			for _, a := range pnmlArcs {
				if a.Source != t.Id && a.Target != t.Id {
					continue
				}

				weight := a.Weight.Int(1)
				arcType := a.Type.String()
				if arcType == "" {
					arcType = a.ArcType.String()
					// CPN has reverted direction of inhibitor edge
					if arcType == "inhibitor" {
						a.Source, a.Target = a.Target, a.Source
					}
				}

				if a.Source == t.Id {
					targets.Push(weight, places.Find(a.Target))
				}
				if a.Target == t.Id {
					if arcType == "inhibitor" {
						origins.PushInhibitor(places.Find(a.Source))
					} else {
						origins.Push(weight, places.Find(a.Source))
					}
				}
			}
			transition := &net.Transition{
				Id:          t.Id,
				Origins:     origins,
				Targets:     targets,
				Priority:    t.Priority,
				Description: t.Name.String(),
				TimeFunc:    nil, // TODO
			}

			composition.Move(transition, t.Position.X, t.Position.Y)
			transitions.Push(transition)
		}
	}

	buildPlaces(pnml.Net.Places)
	buildTransitions(pnml.Net.Transitions, pnml.Net.Arcs)

	for _, page := range pnml.Net.Pages {
		buildPlaces(page.Places)
		buildTransitions(page.Transitions, page.Arcs)
	}

	return net.New(places, transitions), composition
}

func Parse(pnmlReader io.Reader) (net.Net, compose.Composition) {
	pnml := &Pnml{}
	decoder := xml.NewDecoder(pnmlReader)
	decoder.Decode(pnml)
	net, composition := pnml.buildNetCompo()
	composition.CenterTo(0, 0)
	return net, composition
}
