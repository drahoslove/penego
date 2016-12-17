package gui

import (
	"math"
	"image/color"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
)

type OnRedrawFunc func()

var (
	redraw OnRedrawFunc = nil
	ctx * draw2dgl.GraphicContext
)

func init () {
	redraw = func() {
		if ctx != nil {
			ctx.MoveTo(-30, -30)
			ctx.LineTo(+30, +30)
			ctx.Stroke()
			ctx.MoveTo(+30, -30)
			ctx.LineTo(-30, +30)
			ctx.Stroke()
		}
	}
}

var ( // pseudo constants
	WHITE = color.RGBA{255, 255, 255, 255}	// #ffffff
	WHITISH = color.RGBA{239, 239, 239, 255}	// #efefef

	LIGHT_GRAY = color.RGBA{204, 204, 204, 255}	// #cccccc
	GRAY = color.RGBA{128, 128, 128, 255}	// #808080
	DARK_GRAY = color.RGBA{51, 51, 51, 255}	// #333333

	BLACKISH = color.RGBA{16, 16, 16, 255}	// #101010
	BLACK = color.RGBA{0, 0, 0, 255}	// #000000
)

func OnRedraw(f OnRedrawFunc) {
	redraw = f;
	invalid = true
}


func DrawPlace(x, y, n int) {
	if ctx != nil {
		drawPlace(ctx, float64(x), float64(y), n)
	}
}

func DrawTransition(x, y int) {
	if ctx != nil {
		drawTransition(ctx, float64(x), float64(y))
	}
}

func draw() {

	/* create graphic context and set styles */
	ctx = draw2dgl.NewGraphicContext(width, height)
	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.SetLineWidth(3)

	/* background */
	draw2dkit.Rectangle(ctx, 0, 0, float64(width), float64(height))
	ctx.Fill()

	/* translate origin to center */
	ctx.Translate(float64(width/2), float64(height/2))

	// draw shapes
	if redraw != nil {
		redraw()
	}
}

func drawPlace(ctx * draw2dgl.GraphicContext, x float64, y float64, n int) {
	r := 24.0
	ctx.Save()
	defer ctx.Restore()

	ctx.BeginPath()
	draw2dkit.Circle(ctx, x, y, r)
	ctx.Close()

	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.FillStroke()

	if n == 1 {
		draw2dkit.Circle(ctx, x, y, 6)
		ctx.Close()

		ctx.SetFillColor(BLACKISH)
		ctx.Fill()
	}
	if 1 < n && n < 5 {
		rr := float64(r)/(3-float64(n)*0.25)
		for i := 1; i <= n; i++ {
			angle := math.Pi/float64(n)*float64(i)*2
			xx := x+math.Sin(angle)*rr
			yy := y+math.Cos(angle)*rr

			draw2dkit.Circle(ctx, xx, yy, 5)
			ctx.Close()

			ctx.SetFillColor(BLACKISH)
			ctx.Fill()
		}
	}
}

func drawTransition(ctx * draw2dgl.GraphicContext, x float64, y float64) {
	w, h := 18.0, 72.0 // 1:4
	ctx.Save()
	defer ctx.Restore()

	ctx.BeginPath()
	draw2dkit.Rectangle(ctx, x-w/2, y-h/2, x+w/2, y+h/2)
	ctx.Close()
	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.FillStroke()
}