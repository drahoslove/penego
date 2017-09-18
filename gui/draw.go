package gui

// draw content
// drawing routines definitions
// exports nothing

import (
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"math"
	"strconv"
)

const (
	PLACE_RADIUS      = 24.0
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

func opaque (clr color.RGBA, opacity float32) color.RGBA {
	clr.A = uint8(255 * opacity)
	return clr
}

var (
	drawSplash = RedrawFunc(func(screen *Screen) {
		ctx := screen.ctx
		if ctx != nil {
			defer tempContext(ctx)()
			ctx.SetFontData(draw2d.FontData{Name: "gobold"})
			ctx.SetFontSize(48)
			ctx.SetFillColor(DARK_GRAY)
			drawCenteredString(ctx, "Penego", 0, 0)
		}
	})
)

func newCtx(width, height int) *draw2dgl.GraphicContext {
	/* create graphic context and set styles */
	var ctx = draw2dgl.NewGraphicContext(width, height)
	ctx.SetFontData(draw2d.FontData{Name: "goregular"})
	ctx.SetFontSize(14)
	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.SetLineWidth(3)

	/* translate origin to center */
	ctx.Translate(float64(width)/2, float64(height)/2)
	return ctx
}

func clean(ctx *draw2dgl.GraphicContext, width, height int) {
	/* background */
	draw2dkit.Rectangle(ctx, -float64(width)/2, -float64(height)/2, float64(width)/2, float64(height)/2)
	ctx.Fill()
}

func tempContext(ctx *draw2dgl.GraphicContext) (func()) {
	ctx.Save()
	return func () {
		ctx.Restore()
	}
}

// GUI entities

func drawMenu(ctx *draw2dgl.GraphicContext, width, height int, itemsNames []string) {
	defer tempContext(ctx)()
	ctx.Translate(-float64(width)/2, -float64(height)/2)

	ctx.SetFillColor(opaque(DARK_GRAY, 0.9))
	draw2dkit.Rectangle(ctx, 0, 0, float64(width), 36)
	ctx.Fill()

	l := len(itemsNames)
	btnW := width/l

	ctx.SetFillColor(WHITISH)
	ctx.SetFontSize(12)
	for i, name := range itemsNames {
		drawCenteredString(ctx, name, float64(i*btnW + btnW/2), 23)
	}

}


// NET entities

func drawPlace(ctx *draw2dgl.GraphicContext, x float64, y float64, n int, description string) {
	r := PLACE_RADIUS
	defer tempContext(ctx)()

	draw2dkit.Circle(ctx, x, y, r)

	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.FillStroke()

	// tokens
	switch {
	case n == 1: // draw dot
		draw2dkit.Circle(ctx, x, y, 6)
		ctx.SetFillColor(BLACKISH)
		ctx.Fill()
	case 1 < n && n < 6: // draw dots
		rr := float64(r) / (3 - float64(n)*0.25)
		for i := 1; i <= n; i++ {
			angle := math.Pi / float64(n) * float64(i) * 2
			xx := x + math.Sin(angle)*rr
			yy := y + math.Cos(angle)*rr

			draw2dkit.Circle(ctx, xx, yy, 5)
		}
		ctx.SetFillColor(BLACKISH)
		ctx.Fill()

	case n >= 6: // draw numbers
		ctx.Save()
		ctx.SetFontData(draw2d.FontData{Name: "gomono"})
		ctx.SetFillColor(BLACKISH)
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
	defer tempContext(ctx)()

	draw2dkit.Rectangle(ctx, x-w/2, y-h/2, x+w/2, y+h/2)
	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.FillStroke()

	// timed or priority
	if attrs != "" {
		ctx.SetFillColor(BLACKISH)
		drawCenteredString(ctx, attrs, x, y+h/2+20) // under
	}

	// description
	if description != "" {
		ctx.SetFillColor(BLACKISH)
		drawCenteredString(ctx, description, x, y-h/2-10) //ahove
	}
}

func drawArc(ctx *draw2dgl.GraphicContext, fromx, fromy, tox, toy float64, dir Direction, weight int) {
	r := PLACE_RADIUS
	w := TRANSITION_WIDTH
	var cPs []mgl.Vec2 // control point of arcs curve

	defer tempContext(ctx)()

	if dir == In { // ( ) -> [ ]
		angle := math.Pi * +0.25 // outgoing angle from place
		if fromy > toy {
			angle += math.Pi
		}
		xo := math.Sin(angle) * r // start position on place edge related to its center
		yo := math.Cos(angle) * r
		tox -= w / 2
		fromx += xo
		fromy += yo

		cPs = []mgl.Vec2{
			{fromx, fromy},
			{fromx + 4*xo, fromy + 4*yo},
			{tox - 60, toy},
			{tox, toy},
		}
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
		cPs = []mgl.Vec2{
			{fromx, fromy},
			{fromx + 60, fromy},
			{tox + 4*xo, toy + 4*yo},
			{tox, toy},
		}
		drawArrowHead(ctx, tox, toy, angle)
	}

	ctx.MoveTo(cPs[0].X(), cPs[0].Y())
	ctx.CubicCurveTo(
		cPs[1].X(), cPs[1].Y(),
		cPs[2].X(), cPs[2].Y(),
		cPs[3].X(), cPs[3].Y(),
	)
	ctx.Stroke()

	if weight > 1 {
		arcCntr := mgl.CubicBezierCurve2D(0.5, cPs[0], cPs[1], cPs[2], cPs[3])
		draw2dkit.Circle(ctx, arcCntr.X(), arcCntr.Y()-6, 12)
		ctx.SetFillColor(WHITISH)
		ctx.Fill()

		ctx.SetFillColor(BLACKISH)
		drawCenteredString(ctx, strconv.Itoa(weight), arcCntr.X()-1, arcCntr.Y())
	}
}

// help functions

func drawArrowHead(ctx *draw2dgl.GraphicContext, x float64, y float64, angle float64) {
	r := 18.0
	w := math.Pi / 8
	xl := x + math.Sin(angle+w)*r
	yl := y + math.Cos(angle+w)*r
	xr := x + math.Sin(angle-w)*r
	yr := y + math.Cos(angle-w)*r

	defer tempContext(ctx)()

	ctx.MoveTo(x, y)
	ctx.LineTo(xl, yl)
	ctx.LineTo(xr, yr)
	ctx.LineTo(x, y)
	ctx.SetFillColor(BLACKISH)
	ctx.Fill()
}

func drawCenteredString(ctx *draw2dgl.GraphicContext, str string, x float64, y float64) {
	left, top, right, bottom := ctx.GetStringBounds(str)
	width := right - left
	height := bottom - top
	_ = height

	ctx.FillStringAt(str, x-width/2, y)
}
