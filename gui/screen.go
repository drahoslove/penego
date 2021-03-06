package gui

// GUI API functions for drawing items, screen manipulation, event handling, menu init etc.
// exports Screen

import (
	"time"

	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/storage"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/llgcode/draw2d/draw2dgl"
)

var (
	guiSt *storage.Storage
)

func init() {
	guiSt = storage.Of("gui")
}

func nameToKey(key string) glfw.Key {
	switch {
	case key >= "A" && key <= "Z":
		return glfw.Key(rune(key[0])-'A') + glfw.KeyA
	case key == "space":
		return glfw.KeySpace
	case key == "home":
		return glfw.KeyHome
	case key == "right":
		return glfw.KeyRight
	case key == "left":
		return glfw.KeyLeft
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
	menusVisible    bool
	mainMenu        menu
	minorMenu       menu
	hoverTimer      *time.Timer
	style           draw.Style
}

/* non-exported methods */

func (s *Screen) newCtx() {
	s.ctx = draw2dgl.NewGraphicContext(s.width, s.height)
	draw.Init(s.ctx, s.width, s.height)
}

func (s *Screen) drawContent() {
	if s.drawContentFunc != nil {
		draw.Clean(s.ctx, s.width, s.height)
		s.drawContentFunc(s)
		if s.menusVisible {
			draw.ExportBorder(s.ctx)
			menuI := s.mainMenu.activeIndex
			tooltip := s.mainMenu.tooltip()
			widths, height, top := draw.Menu(s.ctx, s.width, s.height,
				s.mainMenu.itemIcons(), menuI, tooltip, s.mainMenu.disabled(), draw.Up)
			s.mainMenu.setBounds(widths, height, top)

			menuI = s.minorMenu.activeIndex
			tooltip = s.minorMenu.tooltip()
			widths, height, top = draw.Menu(s.ctx, s.width, s.height,
				s.minorMenu.itemIcons(), menuI, tooltip, s.minorMenu.disabled(), draw.Down)
			s.minorMenu.setBounds(widths, height, top)
		}
	}
	s.SwapBuffers()
}

func (s *Screen) setActiveMenuIndex(menu *menu, i int) {
	if menu.activeIndex != i {
		menu.activeIndex = i
		menu.showTooltip = false
		s.contentInvalid = true
	}
	if menu.activeIndex == -1 {
		menu.showTooltip = false
		s.contentInvalid = true
	}
	if i != -1 {
		if s.hoverTimer != nil {
			s.hoverTimer.Stop()
		}
		s.hoverTimer = time.AfterFunc(time.Second*3/4, func() {
			menu.showTooltip = true
			s.contentInvalid = true
		})
	}
}

func (s *Screen) setSizeCallback(f func(*Screen, int, int)) {
	s.Window.SetSizeCallback(func(window *glfw.Window, w, h int) {
		f(s, w, h)
	})
}

/* exported methods */

func (s *Screen) Reset() {
	offset := guiSt.Of("offset")
	offset.Set("x", 0.0)
	offset.Set("y", 0.0)
	s.newCtx()
}

func (s *Screen) Pan(dx, dy float64) {
	s.ctx.Translate(dx, dy)
	offset := guiSt.Of("offset")
	offset.AddFloat("x", -dx)
	offset.AddFloat("y", -dy)
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
		s.menusVisible = true
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
		s.menusVisible = false
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

func (s *Screen) SetStyle(style draw.Style) {
	s.style = style
}

func (s *Screen) DrawPlace(pos draw.Pos, n int, description string) {
	if s.ctx != nil {
		draw.Place(s.ctx, s.style, pos, n, description)
	}
}

func (s *Screen) DrawTransition(pos draw.Pos, attrs, description string) {
	if s.ctx != nil {
		draw.Transition(s.ctx, s.style, pos, attrs, description)
	}
}

func (s *Screen) DrawInArc(path []draw.Pos, weight int) {
	if s.ctx != nil {
		draw.Arc(s.ctx, s.style, path, draw.In, weight)
	}
}

func (s *Screen) DrawOutArc(path []draw.Pos, weight int) {
	if s.ctx != nil {
		draw.Arc(s.ctx, s.style, path, draw.Out, weight)
	}
}

func (s *Screen) DrawInhibitorArc(path []draw.Pos) {
	if s.ctx != nil {
		draw.InhibitorArc(s.ctx, s.style, path)
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

func (s *Screen) OnMenu(menu *menu, menuIndex int, cb func()) {
	var prevcb glfw.MouseButtonCallback
	prevcb = s.Window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Release && button == glfw.MouseButton1 {
			if menu.activeIndex == menuIndex {
				doInLoop(cb, false)
			}
		}
		if prevcb != nil {
			prevcb(w, button, action, mod)
		}
	})
}

func (s *Screen) RegisterControl(which int, key string, getIcon func() Icon, label string, handler func(), isEnabled func() bool) {
	menu := []*menu{&s.mainMenu, &s.minorMenu}[which]
	s.OnKey(key, handler)
	i := menu.addItem(getIcon, func() bool { return !isEnabled() }, label)
	s.OnMenu(menu, i, handler)
}

func (s *Screen) normalize(x, y float64, centered bool) (float64, float64) {
	if centered {
		x -= float64(s.width) / 2
		y -= float64(s.height) / 2
		x += guiSt.Float("offset.x")
		y += guiSt.Float("offset.y")
	}
	return x, y
}

func (s *Screen) OnMouseMove(centered bool, cb func(float64, float64) bool) {
	var prevcb glfw.CursorPosCallback
	prevcb = s.Window.SetCursorPosCallback(func(w *glfw.Window, x float64, y float64) {

		if prevcb == nil {
			w.SetCursor(arrowCursor) // set default arrow before first callback
		} else {
			prevcb(w, x, y) // first run callbacks defined earlier
		}
		x, y = s.normalize(x, y, centered)
		if cb(x, y) {
			w.SetCursor(handCursor)
		}
	})
}

func (s *Screen) OnDrag(centered bool, cb func(x, y, deltax, deltaY, startX, startY float64, done bool)) {
	var prevClickCb glfw.MouseButtonCallback
	var prevCurPosCb glfw.CursorPosCallback
	var prevCurEnterCb glfw.CursorEnterCallback

	startX, startY := 0.0, 0.0
	draging := false

	endDrag := func() {
		x, y := s.GetCursorPos()
		x, y = s.normalize(x, y, centered)
		cb(x, y, x-startX, y-startY, startX, startY, true)
		draging = false
	}

	prevClickCb = s.Window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if prevClickCb != nil {
			prevClickCb(w, button, action, mod)
		}
		if button == glfw.MouseButtonLeft {
			if action == glfw.Press {
				draging = true
				startX, startY = w.GetCursorPos()
				startX, startY = s.normalize(startX, startY, centered)
			}
			if action == glfw.Release && draging {
				endDrag()
			}
		}
	})
	prevCurEnterCb = s.Window.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		if prevCurEnterCb != nil {
			prevCurEnterCb(w, entered)
		}
		if !entered && draging {
			endDrag()
		}
	})
	prevCurPosCb = s.Window.SetCursorPosCallback(func(w *glfw.Window, x float64, y float64) {
		if prevCurPosCb != nil {
			prevCurPosCb(w, x, y)
		}
		if draging {
			x, y = s.normalize(x, y, centered)
			cb(x, y, x-startX, y-startY, startX, startY, false)
		}
	})
}
