package gui

// draw content
// drawing routines definitions
// exports nothing

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"math"
	"strconv"
)

const (
	PLACE_READIUS     = 24.0
	TRANSITION_WIDTH  = 18.0
	TRANSITION_HEIGHT = 72.0
)

type RedrawFunc func(*Screen)

var ( // pseudo constants
	WHITE   = color.RGBA{255, 255, 255, 255} // #ffffff
	WHITISH = color.RGBA{239, 239, 239, 255} // #efefef

	LIGHT_GRAY = color.RGBA{204, 204, 204, 255} // #cccccc
	GRAY       = color.RGBA{128, 128, 128, 255} // #808080
	DARK_GRAY  = color.RGBA{51, 51, 51, 255}    // #333333

	BLACKISH = color.RGBA{16, 16, 16, 255} // #101010
	BLACK    = color.RGBA{0, 0, 0, 255}    // #000000
)

var (
	drawSplash = RedrawFunc(func(screen *Screen) {
		ctx := screen.ctx
		if ctx != nil {
			ctx.Save()
			ctx.SetFontData(draw2d.FontData{Name: "gobold"})
			ctx.SetFontSize(48)
			ctx.SetFillColor(DARK_GRAY)
			drawCenteredString(ctx, "Penego", 0, 0)
			ctx.Restore()
		}
	})
)

// called by mainloop if context is invalid
func draw(screen Screen) {

	/* create graphic context and set styles */
	screen.ctx = draw2dgl.NewGraphicContext(screen.width, screen.height)
	ctx := screen.ctx

	ctx.SetFontData(draw2d.FontData{Name: "goregular"})
	ctx.SetFontSize(14)
	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.SetLineWidth(3)

	/* background */
	draw2dkit.Rectangle(ctx, 0, 0, float64(screen.width), float64(screen.height))
	ctx.Fill()

	/* translate origin to center */
	ctx.Translate(float64(screen.width/2), float64(screen.height/2))

	// draw shapes or whatever
	screen.drawContent()
}

func drawPlace(ctx *draw2dgl.GraphicContext, x float64, y float64, n int, description string) {
	r := PLACE_READIUS
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
	case n == 1: // draw dot
		draw2dkit.Circle(ctx, x, y, 6)
		ctx.Close()
		ctx.SetFillColor(BLACKISH)
		ctx.Fill()
	case 1 < n && n < 6: // draw dots
		rr := float64(r) / (3 - float64(n)*0.25)
		for i := 1; i <= n; i++ {
			angle := math.Pi / float64(n) * float64(i) * 2
			xx := x + math.Sin(angle)*rr
			yy := y + math.Cos(angle)*rr

			draw2dkit.Circle(ctx, xx, yy, 5)
			ctx.Close()

			ctx.SetFillColor(BLACKISH)
			ctx.Fill()
		}

	case n >= 6: // draw numbers
		ctx.Save()
		ctx.SetFontData(draw2d.FontData{Name: "gomono"})
		switch {
		case n < 100:
			ctx.SetFontSize(24)
			drawCenteredString(ctx, strconv.Itoa(n), x-1, y+10)
		case n < 1000:
			ctx.SetFontSize(18)
			drawCenteredString(ctx, strconv.Itoa(n), x-1, y+7)
		case n < 10000:
			ctx.SetFontSize(14)
			drawCenteredString(ctx, strconv.Itoa(n), x-1, y+5)
		default:
			// TODO "~XeN" form
			ctx.SetFontSize(14)
			drawCenteredString(ctx, "many", x-1, y+5)
		}
		ctx.Restore()
	}

	// description
	if description != "" {
		drawCenteredString(ctx, description, x, y-r-8)
	}

}

func drawTransition(ctx *draw2dgl.GraphicContext, x, y float64, attrs, description string) {
	w, h := TRANSITION_WIDTH, TRANSITION_HEIGHT
	ctx.Save()
	defer ctx.Restore()

	ctx.BeginPath()
	draw2dkit.Rectangle(ctx, x-w/2, y-h/2, x+w/2, y+h/2)
	ctx.Close()
	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.FillStroke()

	// timed or priority
	if attrs != "" {
		drawCenteredString(ctx, attrs, x, y+h/2+20) // under
	}

	// description
	if description != "" {
		drawCenteredString(ctx, description, x, y-h/2-10) //ahove
	}
}

func drawArc(ctx *draw2dgl.GraphicContext, fromx, fromy, tox, toy float64, dir Direction, weight int) {
	r := PLACE_READIUS
	w := TRANSITION_WIDTH
	if dir == In { // ( ) -> [ ]
		angle := math.Pi * +0.25
		if fromy > toy {
			angle += math.Pi
		}
		xo := math.Sin(angle) * r
		yo := math.Cos(angle) * r
		tox -= w / 2
		fromx += xo
		fromy += yo
		ctx.MoveTo(fromx, fromy)
		ctx.CubicCurveTo(fromx+4*xo, fromy+4*yo, tox-60, toy, tox, toy)
		drawArrowHead(ctx, tox, toy, -math.Pi/2)
	}
	if dir == Out { // [ ] -> ( )
		angle := math.Pi * -0.25
		if fromy < toy {
			angle += math.Pi
		}
		xo := math.Sin(angle) * r
		yo := math.Cos(angle) * r
		fromx += w / 2
		tox += xo
		toy += yo
		ctx.MoveTo(fromx, fromy)
		ctx.CubicCurveTo(fromx+60, fromy, tox+4*xo, toy+4*yo, tox, toy)
		drawArrowHead(ctx, tox, toy, angle)
	}
	ctx.Stroke()
}

// help functions

func drawArrowHead(ctx *draw2dgl.GraphicContext, x float64, y float64, angle float64) {
	r := 18.0
	w := math.Pi / 8
	xl := x + math.Sin(angle+w)*r
	yl := y + math.Cos(angle+w)*r
	xr := x + math.Sin(angle-w)*r
	yr := y + math.Cos(angle-w)*r

	ctx.Save()
	ctx.SetFillColor(BLACKISH)
	ctx.BeginPath()
	ctx.MoveTo(x, y)
	ctx.LineTo(xl, yl)
	ctx.LineTo(xr, yr)
	ctx.LineTo(x, y)
	ctx.Close()
	ctx.Fill()
	ctx.Restore()
}

func drawCenteredString(ctx *draw2dgl.GraphicContext, str string, x float64, y float64) {
	left, top, right, bottom := ctx.GetStringBounds(str)
	width := right - left
	height := bottom - top
	_ = height
	ctx.Save()
	ctx.SetFillColor(BLACKISH)
	ctx.FillStringAt(str, x-width/2, y)
	ctx.Restore()
}
