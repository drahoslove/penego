package gui

// exports GUI API functions for drawing items, screen manipulation, event handling

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/llgcode/draw2d/draw2dgl"
	mgl "github.com/go-gl/mathgl/mgl64"
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
	menuVisible     bool
	menu            Menu
}

// TODO move to its own file
type Menu struct {
	items []MenuItem
	activeIndex int
}

type MenuItem struct {
	name string
	bound Bound
}

type Bound struct {
	from mgl.Vec2 // left top coord
	to mgl.Vec2 // right bottom coord
}

func (b *Bound) hits(x, y float64) bool {
	return x >= b.from.X() && x < b.to.X() &&
		y >= b.from.Y() && y < b.to.Y()
}

func newMenu(names []string) Menu {
	var menu Menu
	menu.items = make([]MenuItem, len(names))
	for i, name := range names {
		menu.items[i] = MenuItem{name:name}
	}
	menu.activeIndex = -1
	return menu
}

func (m *Menu) itemNames() []string {
	var names = make([]string, len(m.items))
	for i, item := range m.items { // TODO this is not sorted
		names[i] = item.name
		i++
	}
	return names
}

func (m *Menu) setBounds(widths []int, height int) {
	from := mgl.Vec2{0, 0}
	to := mgl.Vec2{0, float64(height)}
	for i := range m.items {
		to[0] += float64(widths[i])
		m.items[i].bound = Bound{from, to}
		from[0] += float64(widths[i])
	}
}

func (m *Menu) setActive(index int) {
	m.activeIndex = index
}

func (s *Screen) drawContent() {
	if s.drawContentFunc != nil {
		clean(s.ctx, s.width, s.height)
		s.drawContentFunc(s)
		if s.menuVisible {
			widths, height := drawMenu(s.ctx, s.width, s.height, s.menu.itemNames(), s.menu.activeIndex) // TODO pas button state (select/higlight/pressed)
			s.menu.setBounds(widths, height)
		}
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
		s.menuVisible = true
	}, true)
}

func (s *Screen) SetRedrawFuncToSplash() {
	doInLoop(func() {
		s.drawContentFunc = drawSplash
		s.contentInvalid = true
		s.menuVisible = false
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
		if prevcb != nil {
			prevcb(w, key, scancode, action, mods)
		}
	})
}

func (s *Screen) OnMenu(menuIndex int, cb func()) {
	var prevcb glfw.MouseButtonCallback
	prevcb = s.Window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Release && button == glfw.MouseButton1 {
			if s.menu.activeIndex == menuIndex {
				doInLoop(cb, false)
			}
		}
		if prevcb != nil {
			prevcb(w, button, action, mod)
		}
	})
}