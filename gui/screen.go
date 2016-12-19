package gui

import (
	"github.com/go-gl/glfw/v3.2/glfw"
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

func (s * Screen) DrawPlace(x, y, n int, description string) {
	if ctx != nil {
		drawPlace(ctx, float64(x), float64(y), n, description)
	}
}

func (s * Screen) DrawTransition(x, y int, attrs, description string) {
	if ctx != nil {
		drawTransition(ctx, float64(x), float64(y), attrs, description)
	}
}

func (s * Screen) SetTitle(title string) {
	if title != "" {
		title = " - " + title
	}
	s.Window.SetTitle("Penego" + title)
}