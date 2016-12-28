package gui

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Pos struct {
	X float64
	Y float64
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


func (s * Screen) ForceRedraw(block bool) {
	doInLoop(func() {
		contentInvalid = true
	}, block)
}

func (s * Screen) SetRedrawFunc(f RedrawFunc) {
	doInLoop(func() {
		drawContentFunc = f; // update drawContentFunc
		contentInvalid = true // force draw
	}, true)
}

func (s * Screen) SetTitle(title string) {
	doInLoop(func() {
		if title != "" {
			title = " - " + title
		}
		s.Window.SetTitle("Penego" + title)
	}, false)
}

func (s * Screen) DrawPlace(pos Pos, n int, description string) {
	if ctx != nil {
		drawPlace(ctx, pos.X, pos.Y, n, description)
	}
}

func (s * Screen) DrawTransition(pos Pos, attrs, description string) {
	if ctx != nil {
		drawTransition(ctx, pos.X, pos.Y, attrs, description)
	}
}

func (s * Screen) DrawInArc(from Pos, to Pos, weight int) {
	if ctx != nil {
		drawArc(ctx, from.X, from.Y, to.X, to.Y, In, weight)
	}
}

func (s * Screen) DrawOutArc(from Pos, to Pos, weight int) {
	if ctx != nil {
		drawArc(ctx, from.X, from.Y, to.X, to.Y, Out, weight)
	}
}