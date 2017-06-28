package gui

// exports GUI API functions for drawing items, screen manipulation, event handling

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/llgcode/draw2d/draw2dgl"
)

type Pos struct {
	X float64
	Y float64
}

type Direction bool

const (
	In  Direction = true
	Out Direction = false
)

var (
	nameToKey = map[string]glfw.Key{
		"space": glfw.KeySpace,
		"Q":     glfw.KeyQ,
		"R":     glfw.KeyR,
	}
)

// Screen provide exported functions for drawing graphic content
type Screen struct {
	*glfw.Window
	ctx             *draw2dgl.GraphicContext
	drawContentFunc RedrawFunc
	contentInvalid  bool
	width           int
	height          int
}

func (s *Screen) drawContent() {
	if s.drawContentFunc != nil {
		s.drawContentFunc(s)
	}
}

func (s *Screen) setSizeCallback(f func(*Screen, int, int)) {
	s.Window.SetSizeCallback(func(window *glfw.Window, w, h int) {
		f(s, w, h)
	})
}

func (s *Screen) ForceRedraw(block bool) {
	doInLoop(func() {
		s.contentInvalid = true
	}, block)
}

func (s *Screen) SetRedrawFunc(f RedrawFunc) {
	doInLoop(func() {
		s.drawContentFunc = f   // update drawContentFunc
		s.contentInvalid = true // force draw
	}, true)
}

func (s *Screen) SetRedrawFuncToSplash() {
	doInLoop(func() {
		s.drawContentFunc = drawSplash
		s.contentInvalid = true
	}, true)
}

func (s *Screen) SetTitle(title string) {
	doInLoop(func() {
		if title != "" {
			title = " - " + title
		}
		s.Window.SetTitle("Penego" + title)
	}, false)
}

func (s *Screen) DrawPlace(pos Pos, n int, description string) {
	if s.ctx != nil {
		drawPlace(s.ctx, pos.X, pos.Y, n, description)
	}
}

func (s *Screen) DrawTransition(pos Pos, attrs, description string) {
	if s.ctx != nil {
		drawTransition(s.ctx, pos.X, pos.Y, attrs, description)
	}
}

func (s *Screen) DrawInArc(from Pos, to Pos, weight int) {
	if s.ctx != nil {
		drawArc(s.ctx, from.X, from.Y, to.X, to.Y, In, weight)
	}
}

func (s *Screen) DrawOutArc(from Pos, to Pos, weight int) {
	if s.ctx != nil {
		drawArc(s.ctx, from.X, from.Y, to.X, to.Y, Out, weight)
	}
}

func (s *Screen) OnKey(keyname string, cb func()) {
	var prevcb glfw.KeyCallback
	prevcb = s.Window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press && nameToKey[keyname] == key {
			doInLoop(cb, false)
		}
		prevcb(w, key, scancode, action, mods)
	})
}
