// Package draw contains drawing procedures
// which defines grapical representations of several entities
// including gui elements, such as menu
// and also petri net components - places, transitions, arcs
package draw

import (
	"image/color"
	"math"
	"strconv"

	"git.yo2.cz/drahoslav/penego/storage"
	mgl "github.com/go-gl/mathgl/mgl64"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	settingsSt *storage.Storage
	exportSt   *storage.Storage
	guiSt      *storage.Storage
)

func init() {
	settingsSt = storage.Of("settings")
	exportSt = storage.Of("export")
	guiSt = storage.Of("gui")
}

type Drawer interface {
	DrawPlace(pos Pos, n int, description string)
	DrawTransition(pos Pos, attrs, description string)
	DrawInArc(path []Pos, weight int)
	DrawOutArc(path []Pos, weight int)
	DrawInhibitorArc(path []Pos)
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
	ox, oy := guiSt.Float("offset.x"), guiSt.Float("offset.y")

	ctx.Translate(-ox+float64(width)/2, -oy+float64(height)/2)
}

func Clean(ctx draw2d.GraphicContext, width, height int) {
	defer ctx.SetLineWidth(settingsSt.Float("linewidth"))
	defer tempContext(ctx)()
	ox, oy := guiSt.Float("offset.x"), guiSt.Float("offset.y")
	ctx.Translate(ox, oy)

	w, h := float64(width), float64(height)

	/* background */
	draw2dkit.Rectangle(ctx, -w/2, -h/2, w/2, h/2)
	ctx.Fill()
}

func ExportBorder(ctx draw2d.GraphicContext) {
	defer tempContext(ctx)()
	ox, oy := guiSt.Float("offset.x"), guiSt.Float("offset.y")
	ctx.Translate(ox, oy)

	width, height := float64(exportSt.Int("width")), float64(exportSt.Int("height"))
	draw2dkit.Rectangle(ctx, -width/2, -height/2, width/2, height/2)
	ctx.SetLineWidth(1)
	ctx.SetStrokeColor(DARK_GRAY)
	ctx.Stroke()
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

	ox, oy := guiSt.Float("offset.x"), guiSt.Float("offset.y")
	ctx.Translate(ox-float64(sWidth)/2, oy-float64(sHeight)/2)

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
		offset = -offset
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

func Arc(ctx draw2d.GraphicContext, style Style, path []Pos, dir Direction, weight int) {
	r := PLACE_RADIUS
	w := TRANSITION_WIDTH
	const X, Y = 0, 1

	for i := 0; i < len(path)-1; i++ {
		func() {
			from, to := path[i], path[i+1]
			var cPs []mgl.Vec2 // control point of arcs curve

			defer tempContext(ctx)()

			quad := math.Pi / 2

			if dir == In { // ( ) -> [ ]
				angle := quad // from right by default

				if from.X < to.X { // before
					if from.Y > to.Y {
						angle += quad / 2 // from right-up if place below tran
					}
					if from.Y < to.Y {
						angle -= quad / 2 // from right-down if place above tran
					}
				}
				if from.X > to.X { // after
					angle = -quad // from left
				}

				if from.X == to.X { // in same column
					angle = -quad // from left
					if from.Y > to.Y {
						angle -= quad / 2 // from left-up if place below tran
					}
					if from.Y < to.Y {
						angle += quad / 2 // from left-down if place above tran
					}
				}

				var xo, yo float64
				if i == 0 {
					xo = math.Sin(angle) * r // start position on place edge related to its center
					yo = math.Cos(angle) * r
					from.X += xo
					from.Y += yo
				}
				if i == len(path)-2 {
					to.X -= w / 2
				}
				cPs = []mgl.Vec2{
					{from.X, from.Y},
					{from.X, from.Y},
					{to.X, to.Y},
					{to.X, to.Y},
				}
				if i == 0 {
					cPs[1][X] += 4 * xo
					cPs[1][Y] += 4 * yo
				}
				if i == len(path)-2 {
					cPs[2][X] -= 60
					drawArrowHead(ctx, style, to.X, to.Y, -math.Pi/2)
				}
			}
			if dir == Out { // [ ] -> ( )
				angle := -quad // to left by default

				if from.X < to.X { // after
					if from.Y > to.Y {
						angle += quad / 2 // to left-down if place above tran
					}
					if from.Y < to.Y {
						angle -= quad / 2 // to left-up if place below tran
					}
				}

				if from.X > to.X { // before
					angle = +quad // to right
				}

				if from.X == to.X { // in same column
					angle = +quad
					if from.Y > to.Y {
						angle -= quad / 2 // to right-down if place above tran
					}
					if from.Y < to.Y {
						angle += quad / 2 // to right-up if place below tran
					}
				}
				var xo, yo float64
				if i == 0 {
					from.X += w / 2
				}
				if i == len(path)-2 {
					xo = math.Sin(angle) * r
					yo = math.Cos(angle) * r
					to.X += xo
					to.Y += yo
				}
				cPs = []mgl.Vec2{
					{from.X, from.Y},
					{from.X, from.Y},
					{to.X, to.Y},
					{to.X, to.Y},
				}
				if i == 0 {
					cPs[1][X] += 60
				}
				if i == len(path)-2 {
					cPs[2][X] += 4 * xo
					cPs[2][Y] += 4 * yo
					drawArrowHead(ctx, style, to.X, to.Y, angle)
				}
			}
			if i > 0 { // draw path join
				draw2dkit.Circle(ctx, cPs[0].X(), cPs[0].Y(), 2)
			}

			ctx.MoveTo(cPs[0].X(), cPs[0].Y())
			ctx.CubicCurveTo(
				cPs[1].X(), cPs[1].Y(),
				cPs[2].X(), cPs[2].Y(),
				cPs[3].X(), cPs[3].Y(),
			)
			ctx.SetStrokeColor(style.Color())
			ctx.Stroke()

			if weight > 1 && i == len(path)/2 {
				arcCntr := mgl.CubicBezierCurve2D(0.5, cPs[0], cPs[1], cPs[2], cPs[3])
				draw2dkit.Circle(ctx, arcCntr.X(), arcCntr.Y()-6, 12)
				ctx.SetFillColor(style.Background())
				ctx.Fill()

				ctx.SetFillColor(style.Color())
				drawCenteredString(ctx, strconv.Itoa(weight), arcCntr.X()-1, arcCntr.Y())
			}
		}()
	}
}

func InhibitorArc(ctx draw2d.GraphicContext, style Style, path []Pos) {
	// TODO whole path
	from, to := path[0], path[len(path)-1]
	// Ingibitor edge is alwas from place to transtition: ( ) -> [ ]
	r := PLACE_RADIUS
	w := TRANSITION_WIDTH
	cr := w / 2.8
	var cPs []mgl.Vec2 // control point of arcs curve

	defer tempContext(ctx)()

	angle := math.Pi * +0.25 // from left-up by default
	if from.Y > to.Y {
		angle += math.Pi // from right-down if place above tran
	}
	if from.Y == to.Y && from.X < to.X { // if in line
		angle = math.Pi * 0.5
	}
	xo := math.Sin(angle) * r // start position on place edge related to its center
	yo := math.Cos(angle) * r
	to.X -= w/2 + cr*2
	to.Y += 14 // to not overlap with normal arc arrow
	from.X += xo
	from.Y += yo

	cPs = []mgl.Vec2{
		{from.X, from.Y},
		{from.X + 4*xo, from.Y + 4*yo},
		{to.X - 60, to.Y},
		{to.X, to.Y},
	}

	// drawArrowHead(ctx, style, to.X, to.Y, -math.Pi/2)
	ctx.SetStrokeColor(style.Color())

	// circle
	draw2dkit.Circle(ctx, to.X+cr, to.Y, cr)
	ctx.Stroke()

	ctx.MoveTo(cPs[0].X(), cPs[0].Y())
	ctx.CubicCurveTo(
		cPs[1].X(), cPs[1].Y(),
		cPs[2].X(), cPs[2].Y(),
		cPs[3].X(), cPs[3].Y(),
	)
	ctx.Stroke()

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
