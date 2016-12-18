package gui

import (
	"math"
	"image/color"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
)

type OnRedrawFunc func()

var ( // pseudo constants
	WHITE = color.RGBA{255, 255, 255, 255}	// #ffffff
	WHITISH = color.RGBA{239, 239, 239, 255}	// #efefef

	LIGHT_GRAY = color.RGBA{204, 204, 204, 255}	// #cccccc
	GRAY = color.RGBA{128, 128, 128, 255}	// #808080
	DARK_GRAY = color.RGBA{51, 51, 51, 255}	// #333333

	BLACKISH = color.RGBA{16, 16, 16, 255}	// #101010
	BLACK = color.RGBA{0, 0, 0, 255}	// #000000
)

var (
	drawContent OnRedrawFunc = nil // function for drawing content, settable by OnRedraw
	ctx * draw2dgl.GraphicContext
)

func init () {
	drawContent = func() {
		if ctx != nil {
			// ctx.MoveTo(-30, -30)
			// ctx.LineTo(+30, +30)
			// ctx.Stroke()
			// ctx.MoveTo(+30, -30)
			// ctx.LineTo(-30, +30)
			// ctx.Stroke()
			ctx.Save()
			ctx.SetFontData(draw2d.FontData{Name:"gobold"})
			ctx.SetFontSize(48)
			ctx.SetFillColor(DARK_GRAY)
			drawCenteredString(ctx, "Penego", 0, 0)
			ctx.Restore()
		}
	}
}


func OnRedraw(f OnRedrawFunc) {
	doInLoop(func() {
		drawContent = f; // update drawContentFunc
		contentInvalid = true // force draw
	})
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
	ctx.SetFontData(draw2d.FontData{Name:"goregular"})
	ctx.SetFontSize(16)
	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.SetLineWidth(3)

	/* background */
	draw2dkit.Rectangle(ctx, 0, 0, float64(width), float64(height))
	ctx.Fill()

	/* translate origin to center */
	ctx.Translate(float64(width/2), float64(height/2))

	// draw shapes
	if drawContent != nil {
		drawContent()
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

	// tokens
	switch {
		case n ==1:
			draw2dkit.Circle(ctx, x, y, 6)
			ctx.Close()
			ctx.SetFillColor(BLACKISH)
			ctx.Fill()
		case 1 < n && n < 5:
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
		case n >= 5:

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

func drawCenteredString(ctx * draw2dgl.GraphicContext , str string, x float64, y float64) {
	left, top, right, bottom := ctx.GetStringBounds(str)
	width := right - left
	height := bottom - top
	_ = height
	ctx.FillStringAt(str, x - width/2, y)
}