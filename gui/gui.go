package gui

// exports Run function

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
	"time"
)

var (
	inLoopFuncChan chan func() // TODO remove global
)

func init() {
	runtime.LockOSThread()
	inLoopFuncChan = make(chan func(), 100)
}

func doInLoop(f func(), block bool) {
	done := make(chan bool, 1)
	inLoopFuncChan <- func() {
		f()
		if block {
			done <- true
		}
	}
	if block {
		<-done
	}
}

func Run(handler func(*Screen)) {
	var screen Screen

	// init glfw
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// create window
	displayWidth, displayHeight := getMonitorResolution()
	screen.width, screen.height = displayWidth/2, displayHeight/2
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.Decorated, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)
	window, err := glfw.CreateWindow(screen.width, screen.height, "Penego", nil, nil)
	if err != nil {
		panic(err)
	}
	screen.Window = window
	screen.MakeContextCurrent() // must be called before gl init
	glfw.SwapInterval(1)        // vsync - causes SwapBuffers to wait for frame

	// center window on screen
	screen.SetPos((displayWidth-screen.width)/2, (displayHeight-screen.height)/2)
	screen.Show()

	// init gl
	if err := gl.Init(); err != nil {
		panic(err)
	}

	reshape(&screen, screen.width, screen.height)
	screen.setSizeCallback(reshape)
	screen.SetKeyCallback(onKey)
	screen.SetRefreshCallback(func(window *glfw.Window) {
		screen.drawContent()
		screen.SwapBuffers()
	})

	go func() {
		handler(&screen)
		doInLoop(func() { // close windows after handler returns
			screen.SetShouldClose(true)
		}, false)
	}()

	// main loop
	for !screen.ShouldClose() {
	empty:
		for {
			select {
			case f := <-inLoopFuncChan: // this must be buffer, to not block handler function
				f()
			default:
				break empty
			}
		}

		if screen.contentInvalid {
			screen.drawContent()
			screen.SwapBuffers()
			screen.contentInvalid = false
		}
		glfw.PollEvents()
		time.Sleep(time.Millisecond) // dont waste CPU
	}

}

func reshape(screen *Screen, w, h int) {
	gl.ClearColor(1, 1, 1, 1) // white
	/* Establish viewing area to cover entire window. */
	gl.Viewport(0, 0, int32(w), int32(h))
	/* PROJECTION Matrix mode. */
	gl.MatrixMode(gl.PROJECTION)
	/* Reset project matrix. */
	gl.LoadIdentity()
	/* Map abstract coords directly to window coords. */
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	/* Invert Y axis so increasing Y goes down. */
	gl.Scalef(1, -1, 1)
	/* Shift origin */
	gl.Translatef(0, float32(-h), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)

	screen.width, screen.height = w, h
	screen.ctx = newCtx(w, h)
	screen.contentInvalid = true

}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch {
		case key == glfw.KeyEscape, key == glfw.KeyQ:
			w.SetShouldClose(true)
		}
	}
}

func getMonitorResolution() (int, int) {
	monitor := glfw.GetPrimaryMonitor()
	if monitor == nil {
		return 800, 600
	}
	vidMode := monitor.GetVideoMode()
	if vidMode == nil {
		return 800, 600
	}
	return vidMode.Width, vidMode.Height
}
