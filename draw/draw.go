package draw

// draw content
// drawing routines definitions
// exports:
//   Drawer
//   Pos, Direction
//   Init, Clean, Splash, Menu
//   Place, Transition, Arc

import (
	"git.yo2.cz/drahoslav/penego/storage"
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
	"image/color"
	"math"
	"strconv"
)

var (
	settingsSt *storage.Storage
)

func init() {
	settingsSt = storage.Of("settings")
}

type Drawer interface {
	DrawPlace(pos Pos, n int, description string)
	DrawTransition(pos Pos, attrs, description string)
	DrawInArc(from, to Pos, weight int)
	DrawOutArc(from, to Pos, weight int)
	SetStyle(style Style)
}

type Style int

const (
	DefaultStyle = Style(iota)
	HighlightedStyle
	FadedStyle
)

func (s Style) Color() color.RGBA {
	return map[Style]color.RGBA{
		DefaultStyle:     BLACKISH,
		HighlightedStyle: BLACK,
		FadedStyle:       opaque(DARK_GRAY, 0.5),
	}[s]
}

func (s Style) Background() color.RGBA {
	return map[Style]color.RGBA{
		DefaultStyle:     WHITISH,
		HighlightedStyle: WHITISH,
		FadedStyle:       opaque(WHITISH, 0.0),
	}[s]
}

type Pos struct {
	X float64
	Y float64
}

func (pos Pos) Equal(p Pos) bool {
	return pos.X == p.X && pos.Y == p.Y
}

type Gravity bool

const (
	Up   Gravity = true
	Down Gravity = false
)

type Direction bool

const (
	In  Direction = true
	Out Direction = false
)

const (
	PLACE_RADIUS      = 24.0
	TRANSITION_WIDTH  = 18.0
	TRANSITION_HEIGHT = 72.0
)

var ( // pseudo constants
	WHITE   = color.RGBA{255, 255, 255, 255} // #ffffff
	WHITISH = color.RGBA{239, 239, 239, 255} // #efefef

	LIGHT_GRAY = color.RGBA{204, 204, 204, 255} // #cccccc
	GRAY       = color.RGBA{128, 128, 128, 255} // #808080
	DARK_GRAY  = color.RGBA{51, 51, 51, 255}    // #333333

	BLACKISH = color.RGBA{16, 16, 16, 255} // #101010
	BLACK    = color.RGBA{0, 0, 0, 255}    // #000000
)

func Init(ctx draw2d.GraphicContext, width, height int) {
	/* create graphic context and set styles */
	ctx.SetFontData(draw2d.FontData{Name: "goregular"})
	ctx.SetFontSize(14)
	ctx.SetFillColor(WHITISH)
	ctx.SetStrokeColor(BLACKISH)
	ctx.SetLineWidth(settingsSt.Float("linewidth"))

	/* translate origin to center */
	ctx.Translate(float64(width)/2, float64(height)/2)
}

func Clean(ctx draw2d.GraphicContext, width, height int) {
	/* background */
	draw2dkit.Rectangle(ctx, -float64(width)/2, -float64(height)/2, float64(width)/2, float64(height)/2)
	ctx.Fill()
	ctx.SetLineWidth(settingsSt.Float("linewidth"))
}

// GUI entities

func Splash(ctx draw2d.GraphicContext, title string) {
	defer tempContext(ctx)()
	ctx.SetFontData(draw2d.FontData{Name: "gobold"})
	ctx.SetFontSize(48)
	ctx.SetFillColor(DARK_GRAY)
	drawCenteredString(ctx, title, 0, 0)
}

func Menu(ctx draw2d.GraphicContext, sWidth, sHeight int, itemsNames []string, activeIndex int, tooltip string, disabled []bool, pos Gravity) ([]float64, float64, float64) {
	defer tempContext(ctx)()

	const (
		padding = 16.0
		height  = 42.0
	)

	var widths = make([]float64, len(itemsNames))
	top := 0.0
	bot := height
	if pos == Down {
		top = float64(sHeight) - height
		bot = float64(sHeight)
	}

	ctx.Translate(-float64(sWidth)/2, -float64(sHeight)/2)

	ctx.SetFillColor(opaque(DARK_GRAY, 0.9))

	draw2dkit.Rectangle(ctx, 0, top, float64(sWidth), bot)

	ctx.Fill()

	ctx.SetFontSize(18) // 24px
	ctx.SetFontData(draw2d.FontData{Name: "ico"})
	sum := 0.0
	for i, name := range itemsNames {
		if i == activeIndex {
			ctx.SetFillColor(WHITE)
		} else {
			ctx.SetFillColor(LIGHT_GRAY)
		}
		if disabled[i] {
			ctx.SetFillColor(DARK_GRAY)
		}
		linePos := top + height*5/7
		w := ctx.FillStringAt(name, sum+padding, linePos)
		width := math.Ceil(w) + 2*padding
		widths[i] = width
		if tooltip != "" && i == activeIndex {
			drawTooltip(ctx, tooltip, pos, sum, top, width, height)
		}
		sum += width
	}
	return widths, height, top
}

func drawTooltip(ctx draw2d.GraphicContext, text string, pos Gravity, left, top, width, height float64) {
	defer tempContext(ctx)()
	offset := 7.0

	if pos == Up {
		top += height + offset
	}
	if pos == Down {
		top -= height/2 + offset
		offset = -offset
	}
	height /= 2

	ctx.SetFontData(draw2d.FontData{Name: "goregular"})
	ctx.SetFontSize(12)

	_, _, w, _ := ctx.GetStringBounds(text)
	center := left + width/2
	left = center - w/2 - 5
	width = w + 10

	ctx.SetFillColor(GRAY)
	draw2dkit.Rectangle(ctx, left, top, left+width, top+height)
	ctx.Fill()
	ctx.SetFillColor(WHITISH)
	ctx.FillStringAt(text, left+5, top+height-6)

	// rectangle
	if pos == Down {
		top += height
	}
	ctx.MoveTo(center, top-offset)
	ctx.LineTo(center+offset, top)
	ctx.LineTo(center-offset, top)
	ctx.Close()
	ctx.SetFillColor(GRAY)
	ctx.Fill()
}

// NET entities

func Place(ctx draw2d.GraphicContext, style Style, pos Pos, n int, description string) {
	r := PLACE_RADIUS
	x, y := pos.X, pos.Y
	defer tempContext(ctx)()

	draw2dkit.Circle(ctx, x, y, r)

	ctx.SetFillColor(style.Background())
	ctx.SetStrokeColor(style.Color())
	ctx.FillStroke()

	// tokens
	switch {
	case n == 1: // draw dot
		draw2dkit.Circle(ctx, x, y, 6)
		ctx.SetFillColor(style.Color())
		ctx.Fill()
	case 1 < n && n < 6: // draw dots
		rr := float64(r) / (3 - float64(n)*0.25)
		for i := 1; i <= n; i++ {
			angle := math.Pi / float64(n) * float64(i) * 2
			xx := x + math.Sin(angle)*rr
			yy := y + math.Cos(angle)*rr

			draw2dkit.Circle(ctx, xx, yy, 5)
		}
		ctx.SetFillColor(style.Color())
		ctx.Fill()

	case n >= 6: // draw numbers
		ctx.Save()
		ctx.SetFontData(draw2d.FontData{Name: "gomono"})
		ctx.SetFillColor(style.Color())
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
		ctx.SetFillColor(style.Color())
		drawCenteredString(ctx, description, x, y-r-8)
	}

}

func Transition(ctx draw2d.GraphicContext, style Style, pos Pos, attrs, description string) {
	w, h := TRANSITION_WIDTH, TRANSITION_HEIGHT
	x, y := pos.X, pos.Y
	defer tempContext(ctx)()

	draw2dkit.Rectangle(ctx, x-w/2, y-h/2, x+w/2, y+h/2)
	ctx.SetFillColor(style.Background())
	ctx.SetStrokeColor(style.Color())
	ctx.FillStroke()

	// timed or priority
	if attrs != "" {
		ctx.SetFillColor(style.Color())
		drawCenteredString(ctx, attrs, x, y+h/2+20) // under
	}

	// description
	if description != "" {
		ctx.SetFillColor(style.Color())
		drawCenteredString(ctx, description, x, y-h/2-10) //ahove
	}
}

func Arc(ctx draw2d.GraphicContext, style Style, from, to Pos, dir Direction, weight int) {
	r := PLACE_RADIUS
	w := TRANSITION_WIDTH
	var cPs []mgl.Vec2 // control point of arcs curve

	defer tempContext(ctx)()

	if dir == In { // ( ) -> [ ]
		angle := math.Pi * +0.25 // from left-up by default
		if from.Y > to.Y {
			angle += math.Pi // from right-down if place above tran
		}
		if from.Y == to.Y && from.X < to.X { // if in line
			angle = math.Pi * 0.5
		}
		xo := math.Sin(angle) * r // start position on place edge related to its center
		yo := math.Cos(angle) * r
		to.X -= w / 2
		from.X += xo
		from.Y += yo

		cPs = []mgl.Vec2{
			{from.X, from.Y},
			{from.X + 4*xo, from.Y + 4*yo},
			{to.X - 60, to.Y},
			{to.X, to.Y},
		}
		drawArrowHead(ctx, style, to.X, to.Y, -math.Pi/2)
	}
	if dir == Out { // [ ] -> ( )
		angle := math.Pi * -0.25 // to left-down by default
		if from.Y < to.Y {
			angle += math.Pi // to fith-up if place above tran
		}
		if from.Y == to.Y && from.X < to.X { // if in line
			angle = math.Pi * -0.5
		}
		xo := math.Sin(angle) * r
		yo := math.Cos(angle) * r
		from.X += w / 2
		to.X += xo
		to.Y += yo
		cPs = []mgl.Vec2{
			{from.X, from.Y},
			{from.X + 60, from.Y},
			{to.X + 4*xo, to.Y + 4*yo},
			{to.X, to.Y},
		}
		drawArrowHead(ctx, style, to.X, to.Y, angle)
	}

	ctx.MoveTo(cPs[0].X(), cPs[0].Y())
	ctx.CubicCurveTo(
		cPs[1].X(), cPs[1].Y(),
		cPs[2].X(), cPs[2].Y(),
		cPs[3].X(), cPs[3].Y(),
	)
	ctx.SetStrokeColor(style.Color())
	ctx.Stroke()

	if weight > 1 {
		arcCntr := mgl.CubicBezierCurve2D(0.5, cPs[0], cPs[1], cPs[2], cPs[3])
		draw2dkit.Circle(ctx, arcCntr.X(), arcCntr.Y()-6, 12)
		ctx.SetFillColor(style.Background())
		ctx.Fill()

		ctx.SetFillColor(style.Color())
		drawCenteredString(ctx, strconv.Itoa(weight), arcCntr.X()-1, arcCntr.Y())
	}
}

// help functions

func drawArrowHead(ctx draw2d.GraphicContext, style Style, x, y float64, angle float64) {
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
	ctx.SetFillColor(style.Color())
	ctx.Fill()
}

func drawCenteredString(ctx draw2d.GraphicContext, str string, x, y float64) {
	left, top, right, bottom := ctx.GetStringBounds(str)
	width := right - left
	height := bottom - top
	_ = height

	ctx.FillStringAt(str, x-width/2, y)
}

func tempContext(ctx draw2d.GraphicContext) func() {
	ctx.Save()
	return ctx.Restore
}

func opaque(clr color.RGBA, opacity float32) color.RGBA {
	clr.A = uint8(255 * opacity)
	return clr
}
