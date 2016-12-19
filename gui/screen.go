package gui

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Pos struct {
	X int
	Y int
}

type Direction bool

const (
	In Direction = true
	Out Direction = false
)

// Screen provide exported functions for drawing graphic content
type Screen struct{
	*glfw.Window
}


func (s * Screen) ForceRedraw() {
	doInLoop(func() {
		contentInvalid = true
	})
}

func (s * Screen) SetRedrawFunc(f RedrawFunc) {
	doInLoop(func() {
		drawContentFunc = f; // update drawContentFunc
		contentInvalid = true // force draw
	})
}

func (s * Screen) SetTitle(title string) {
	if title != "" {
		title = " - " + title
	}
	s.Window.SetTitle("Penego" + title)
}

func (s * Screen) DrawPlace(pos Pos, n int, description string) {
	if ctx != nil {
		drawPlace(ctx, float64(pos.X), float64(pos.Y), n, description)
	}
}

func (s * Screen) DrawTransition(pos Pos, attrs, description string) {
	if ctx != nil {
		drawTransition(ctx, float64(pos.X), float64(pos.Y), attrs, description)
	}
}

func (s * Screen) DrawInArc(from Pos, to Pos, weight int) {
	if ctx != nil {
		drawArc(ctx, float64(from.X), float64(from.Y), float64(to.X), float64(to.Y), In, weight)
	}
}

func (s * Screen) DrawOutArc(from Pos, to Pos, weight int) {
	if ctx != nil {
		drawArc(ctx, float64(from.X), float64(from.Y), float64(to.X), float64(to.Y), Out, weight)
	}
}