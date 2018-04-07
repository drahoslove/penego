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

// NOTE
//   PIPE uses <value>
//   CPN uses <text>

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
	Id         string   `xml:"id,attr"`
	NameTxt    string   `xml:"name>text,omitempty"`
	NameVal    string   `xml:"name>value,omitempty"`
	MarkingTxt int      `xml:"initialMarking>text"`
	MarkingVal string   `xml:"initialMarking>value"`
	Position   Position `xml:"graphics>position"`
}

type Position struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}

type Transition struct {
	Id       string   `xml:"id,attr"`
	NameTxt  string   `xml:"name>text"`
	NameVal  string   `xml:"name>value"`
	Priority int      `xml:"priority>value"`
	Position Position `xml:"graphics>position"`
}

type Arc struct {
	Id        string `xml:"id,attr"`
	Source    string `xml:"source,attr"`
	Target    string `xml:"target,attr"`
	WeightTxt int    `xml:"inscription>text"`
	WeightVal string `xml:"inscription>value"`
}

func (pnml *Pnml) buildNetCompo() (net.Net, compose.Composition) {
	composition := compose.New()

	places := net.Places{}
	transitions := net.Transitions{}

	buildPlaces := func(pnmlPlaces []Place) {
		for _, p := range pnmlPlaces {
			place := &net.Place{
				Tokens:      p.MarkingTxt,
				Id:          p.Id,
				Description: p.NameTxt,
			}
			if p.NameVal != "" {
				place.Description = p.NameVal
			}
			if p.MarkingVal != "" {
				parts := strings.SplitAfter(p.MarkingVal, "Default,")
				tokens, _ := strconv.Atoi(parts[len(parts)-1])
				place.Tokens = tokens
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
				weight := 1
				if a.WeightTxt > 0 {
					weight = a.WeightTxt
				}
				if a.WeightVal != "" {
					parts := strings.SplitAfter(a.WeightVal, "Default,")
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
				Description: t.NameTxt,
				TimeFunc:    nil, // TODO
			}
			if t.NameVal != "" {
				transition.Description = t.NameVal
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
