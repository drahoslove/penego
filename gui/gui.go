package gui
// exports Run function


import (
	"time"
	"runtime"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	width = 800
	height = 600
	contentInvalid = false
	inLoopFuncChan chan func()
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
	var window * glfw.Window


	// init glfw
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

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
	glfw.SwapInterval(1) // vsync - causes SwapBuffers to wait for frame

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

	go func() {
		handler(&Screen{window})
		doInLoop(func() { // close windows after handler returns
			window.SetShouldClose(true)
		}, false)
	}()

	// main loop
	for !window.ShouldClose() {
		empty: for {
			select {
			case f := <-inLoopFuncChan: // this must be buffer, to not block handler function
				f()
			default:
				break empty
			}
		}

		if contentInvalid {
			draw()
			window.SwapBuffers()
			contentInvalid = false
		}
		glfw.PollEvents()
		time.Sleep(time.Millisecond) // dont waste CPU
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
	gl.Scalef(1, -1, 1)
	/* Shift origin */
	gl.Translatef(0, float32(-h), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)

	width, height = w, h
	contentInvalid = true
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

