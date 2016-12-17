package gui

import (
	"runtime"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)


var (
	width = 800
	height = 600
	invalid = true
)


func init() {
	runtime.LockOSThread()
}

func Run() {

	// init glfw
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// glfw.SwapInterval(1) // vsync

	// create window
	screenWidth, screenHeight := getResolution()
	width, height = screenWidth/2, screenHeight/2
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.Decorated, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)
	window, err := glfw.CreateWindow(width, height, "Penego", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent() // must be called before gl init

	// center window on screen
	window.SetPos((screenWidth-width)/2, (screenHeight-height)/2)
	window.Show()

	// init gl
	if err := gl.Init(); err != nil {
		panic(err)
	}

	reshape(window, width, height)
	window.SetSizeCallback(reshape)
	window.SetKeyCallback(onKey)
	window.SetRefreshCallback(func (window * glfw.Window) {
		draw()
		window.SwapBuffers()
	});

	// main loop
	for !window.ShouldClose() {
		if invalid {
			draw()
			window.SwapBuffers()
			invalid = false
		}
		glfw.PollEvents()
	}
}

func reshape(window *glfw.Window, w, h int) {
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
	// gl.Scalef(1, -1, 1)
	/* Shift origin */
	// gl.Translatef(float32(w/2), float32(h/2), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)

	width, height = w, h
	invalid = true
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch {
		case key == glfw.KeyEscape,	key == glfw.KeyQ:
			w.SetShouldClose(true)
		}
	}
}


func getResolution() (int, int) {
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

