package main

import (
	"fmt"
	"time"
	"penego/net"

	"runtime"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
)



func main() {

	var (
		network net.Net
		err error
	)

	// this petri net:

	/**
	 *
	 *   (1)<-----
	 *    |       |
	 *    |       |    2    exit
	 *     ----->[ ]------->( )
	 *         exp(30s)
	 */

	if true {
		// can be done likek this:

		g := &net.Place{Tokens:1} // generator
		e := &net.Place{Description: "exit"}
		t := &net.Transition{
			Origins: net.Arcs{{1,g}},
			Targets: net.Arcs{{1,g},{2,e}},
			TimeFunc: net.GetExponentialTimeFunc(30*time.Second),
		}
		network = net.New(net.Places{g, e}, net.Transitions{t})
	} else {
		// or like this:
		network, err = net.Parse(`
			g (1)
			e ( ) "exit"
			----
			g -> [exp(30us)] -> g, 2*e
		`)
		if err != nil {
			panic(err)
		}
	}



	////////////////////////////////

	fmt.Println(network)

	sim := net.NewSimulation(0, time.Millisecond, network)
	sim.DoEveryTime = func () {
		fmt.Println(sim.GetNow(), network.Places())
	}

	for i := 0; i < 10; i++ {
		net.TrueRandomSeed()
		sim.Run()
	}

	////////////////////////////////
	gui()

}

var (
	width = 800
	height = 600
)

func init() {
	runtime.LockOSThread()
}

func gui() {

	// init glfw
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// create window
	glfw.WindowHint(glfw.Resizable, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Penego", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1) // vsync

	// init gl
	err = gl.Init()
	if err != nil {
		panic(err)
	}

	reshape(window, width, height)
	window.SetSizeCallback(reshape)
	window.SetKeyCallback(onKey)
	window.SetRefreshCallback(func (window * glfw.Window) {
		display()
		window.SwapBuffers()
	});


	// loop

	for !window.ShouldClose() {
		display()
		window.SwapBuffers()
		glfw.PollEvents()
		// time.Sleep(1 * time.Second/60)
	}

}

func display() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	ctx := draw2dgl.NewGraphicContext(width, height)
	ctx.SetFillColor(color.RGBA{255, 255, 255, 255})
	ctx.SetStrokeColor(color.RGBA{8, 8, 8, 255})

	ctx.SetMatrixTransform(draw2d.NewTranslationMatrix(float64(width/2), float64(height/2)))
	ctx.SetLineWidth(4)

	ctx.BeginPath()
	ctx.MoveTo(-50, 20)
	ctx.LineTo(40, 60)
	// ctx.MoveTo(40, 60)
	ctx.LineTo(75, 0)
	// ctx.MoveTo(75, 60)
	ctx.LineTo(76, 460)
	ctx.Close()
	ctx.Stroke()

	// ctx.BeginPath()
	draw2dkit.Circle(ctx, 0, 0, 24)
	// ctx.Close()
	ctx.Stroke()

	// ctx.BeginPath()
	draw2dkit.Circle(ctx, 50, -50, 24)
	// ctx.Close()

	ctx.Stroke()

	gl.Flush() /* Single buffered, so needs a flush. */
}

func reshape(window *glfw.Window, w, h int) {
	gl.ClearColor(1, 1, 1, 1)
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
	/* Shift origin up to upper-left corner. */
	// gl.Translatef(float32(w/2), float32(h/2), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)

	width, height = w, h
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	}
}