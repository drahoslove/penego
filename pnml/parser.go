package pnml

import (
	"encoding/xml"
	"io"
	"strconv"
	"strings"

	"git.yo2.cz/drahoslav/penego/net"
	"git.yo2.cz/drahoslav/penego/compose"
)

// structures defining pnml format

type Pnml struct {
	Net Net `xml:"net"`
}

type Net struct {
	Places      []Place      `xml:"place"`
	Transitions []Transition `xml:"transition"`
	Arcs        []Arc        `xml:"arc"`
}

type Place struct {
	Id          string `xml:"id,attr"`
	Name        string `xml:"name>value,omitempty"`
	Marking     int    `xml:"initialMarking>text"`
	MarkingPIPE string `xml:"initialMarking>value"`
	Position	Position `xml:"graphics>position"`
}

type Position struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}

type Transition struct {
	Id       string `xml:"id,attr"`
	Name     string `xml:"name>value"`
	Priority int    `xml:"priority>value"` // PIPE
	Position	Position `xml:"graphics>position"`
}

type Arc struct {
	Id         string `xml:"id,attr"`
	Source     string `xml:"source,attr"`
	Target     string `xml:"target,attr"`
	Weight     int    `xml:"inscription>text"`
	WeightPIPE string `xml:"inscription>value"`
}

func (pnml *Pnml) buildNetCompo() (net.Net, compose.Composition) {
	composition := compose.New()

	places := net.Places{}
	transitions := net.Transitions{}

	for _, p := range pnml.Net.Places {
		place := &net.Place{
			Tokens:      p.Marking,
			Id:          p.Id,
			Description: p.Name,
		}
		if p.MarkingPIPE != "" {
			parts := strings.SplitAfter(p.MarkingPIPE, "Default,")
			if len(parts) == 2 {
				tokens, _ := strconv.Atoi(parts[1])
				place.Tokens = tokens
			}
		}

		composition.Move(place, p.Position.X, p.Position.Y)
		places.Push(place)
	}
	for _, t := range pnml.Net.Transitions {
		origins := net.Arcs{}
		targets := net.Arcs{}

		for _, a := range pnml.Net.Arcs {
			weight := 1
			if a.Weight > 0 {
				weight = a.Weight
			}
			if a.WeightPIPE != "" {
				parts := strings.SplitAfter(a.WeightPIPE, "Default,")
				if len(parts) == 2 {
					weight, _ = strconv.Atoi(parts[1])
				}
			}
			if a.Source == t.Id {
				targets.Push(weight, places.Find(a.Target))
			}
			if a.Target == t.Id {
				origins.Push(weight, places.Find(a.Source))
			}
		}
		transition := &net.Transition{
			Id:          t.Id,
			Origins:     origins,
			Targets:     targets,
			Priority:    t.Priority,
			Description: t.Name,
			TimeFunc:    nil, // TODO
		}

		composition.Move(transition, t.Position.X, t.Position.Y)
		transitions.Push(transition)
	}
	return net.New(places, transitions), composition
}

func Parse(pnmlReader io.Reader) (net.Net, compose.Composition) {
	pnml := &Pnml{}
	decoder := xml.NewDecoder(pnmlReader)
	decoder.Decode(pnml)
	net, composition := pnml.buildNetCompo()
	return net, composition
}
