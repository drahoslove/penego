package pnml

import (
	"encoding/xml"
	"fmt"
	"git.yo2.cz/drahoslav/penego/net"
	"io"
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
	Id      string `xml:"id,attr"`
	Name    string `xml:"name>value,omitempty"`
	Marking int    `xml:"initialMarking>text"`
}

type Transition struct {
	Id string `xml:"id,attr"`
}

type Arc struct {
	Id     string `xml:"id,attr"`
	Source string `xml:"source,attr"`
	Target string `xml:"target,attr"`
	Weight int    `xml:"inscription>text"`
}

func (pnml *Pnml) build() *net.Net {
	places := net.Places{}
	transitions := net.Transitions{}
	_ = transitions

	for _, p := range pnml.Net.Places {
		places.Push(&net.Place{
			Tokens:      p.Marking,
			Id:          p.Id,
			Description: p.Name,
		})
	}
	for _, t := range pnml.Net.Transitions {
		origins := net.Arcs{}
		targets := net.Arcs{}

		for _, a := range pnml.Net.Arcs {
			if a.Source == t.Id {
				targets.Push(a.Weight, places.Find(a.Target))
			}
			if a.Target == t.Id {
				origins.Push(a.Weight, places.Find(a.Source))
			}
		}
		transitions.Push(&net.Transition{})
	}
	net := net.New(places, transitions)
	return &net
}

func Parse(pnmlReader io.Reader) *net.Net {
	pnml := &Pnml{}
	decoder := xml.NewDecoder(pnmlReader)
	decoder.Decode(pnml)
	fmt.Printf("%+v", pnml)
	return pnml.build()
}
