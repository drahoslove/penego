package gui

// GUI API functions for drawing items, screen manipulation, event handling, menu init etc.
// exports Screen

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/llgcode/draw2d/draw2dgl"
	"git.yo2.cz/drahoslav/penego/draw"
)

func nameToKey(key string) glfw.Key {
	switch {
	case key >= "A" && key <= "Z":
		return glfw.Key(rune(key[0])-'A') + glfw.KeyA
	case key == "space":
		return glfw.KeySpace
	default:
		return glfw.KeyUnknown
	}
}

type RedrawFunc func(draw.Drawer)

// Screen provide exported functions for drawing graphic content
type Screen struct {
	*glfw.Window
	ctx             *draw2dgl.GraphicContext
	drawContentFunc RedrawFunc
	contentInvalid  bool
	width           int
	height          int
	menuVisible     bool
	menu            menu
}

/* non-exported methods */

func (s *Screen) newCtx(w, h int) {
	s.ctx = draw2dgl.NewGraphicContext(w, h)
	draw.Init(s.ctx, w, h)
}

func (s *Screen) drawContent() {
	if s.drawContentFunc != nil {
		draw.Clean(s.ctx, s.width, s.height)
		s.drawContentFunc(s)
		if s.menuVisible {
			widths, height := draw.Menu(s.ctx, s.width, s.height, s.menu.itemIcons(), s.menu.activeIndex)
			s.menu.setBounds(widths, height)
		}
	}
}

func (s *Screen) setActiveMenuIndex(i int) {
	if s.menu.activeIndex != i {
		s.menu.activeIndex = i
		s.contentInvalid = true
	}
}

func (s *Screen) setSizeCallback(f func(*Screen, int, int)) {
	s.Window.SetSizeCallback(func(window *glfw.Window, w, h int) {
		f(s, w, h)
	})
}

/* exported methods */

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

func (s *Screen) SetRedrawFuncToSplash(title string) {
	doInLoop(func() {
		s.drawContentFunc = RedrawFunc(func(drawer draw.Drawer) {
			ctx := s.ctx
			if ctx != nil {
				draw.Splash(ctx, title)
			}
		})
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

func (s *Screen) DrawPlace(pos draw.Pos, n int, description string) {
	if s.ctx != nil {
		draw.Place(s.ctx, pos, n, description)
	}
}

func (s *Screen) DrawTransition(pos draw.Pos, attrs, description string) {
	if s.ctx != nil {
		draw.Transition(s.ctx, pos, attrs, description)
	}
}

func (s *Screen) DrawInArc(from draw.Pos, to draw.Pos, weight int) {
	if s.ctx != nil {
		draw.Arc(s.ctx, from, to, draw.In, weight)
	}
}

func (s *Screen) DrawOutArc(from draw.Pos, to draw.Pos, weight int) {
	if s.ctx != nil {
		draw.Arc(s.ctx, from, to, draw.Out, weight)
	}
}

func (s *Screen) OnKey(keyName string, cb func()) {
	var prevcb glfw.KeyCallback
	prevcb = s.Window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press && nameToKey(keyName) == key {
			doInLoop(cb, false)
		}
		if prevcb != nil {
			prevcb(w, key, scanCode, action, mods)
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

func (s *Screen) RegisterControl(key string, getIcon func() Icon, label string, handler func()) {
	s.OnKey(key, handler)
	i := s.menu.addItem(getIcon, label)
	s.OnMenu(i, handler)
}
