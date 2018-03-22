package pnml

import (
	"encoding/xml"
	"fmt"
	"git.yo2.cz/drahoslav/penego/net"
	"io"
	"strconv"
	"strings"
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
}

type Transition struct {
	Id       string `xml:"id,attr"`
	Name     string `xml:"name>value"`
	Priority int    `xml:"priority>value"` // PIPE
}

type Arc struct {
	Id         string `xml:"id,attr"`
	Source     string `xml:"source,attr"`
	Target     string `xml:"target,attr"`
	Weight     int    `xml:"inscription>text"`
	WeightPIPE string `xml:"inscription>value"`
}

func (pnml *Pnml) build() *net.Net {
	places := net.Places{}
	transitions := net.Transitions{}
	_ = transitions

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
		places.Push(place)
	}
	for _, t := range pnml.Net.Transitions {
		origins := net.Arcs{}
		targets := net.Arcs{}

		for _, a := range pnml.Net.Arcs {
			weight := a.Weight
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
		transitions.Push(&net.Transition{
			Origins:     origins,
			Targets:     targets,
			Priority:    t.Priority,
			Description: t.Name,
			TimeFunc:    nil, // TODO
		})
	}
	net := net.New(places, transitions)
	return &net
}

func Parse(pnmlReader io.Reader) *net.Net {
	pnml := &Pnml{}
	decoder := xml.NewDecoder(pnmlReader)
	decoder.Decode(pnml)
	fmt.Printf("%+v", pnml.Net.Arcs)
	return pnml.build()
}
