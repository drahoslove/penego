package export

import (
	"image"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"git.yo2.cz/drahoslav/penego/compose"
	"git.yo2.cz/drahoslav/penego/draw"
)

var (
	filename = "img.png"
	width, height = 1024, 512
)

type ImgDrawer struct {
	ctx draw2d.GraphicContext
}

func (drawer ImgDrawer) DrawPlace(pos draw.Pos, n int, description string) {
	if drawer.ctx != nil {
		draw.Place(drawer.ctx, pos, n, description)
	}
}

func (drawer ImgDrawer) DrawTransition(pos draw.Pos, attrs, description string) {
	if drawer.ctx != nil {
		draw.Transition(drawer.ctx, pos, attrs, description)
	}
}

func (drawer ImgDrawer) DrawInArc(from draw.Pos, to draw.Pos, weight int) {
	if drawer.ctx != nil {
		draw.Arc(drawer.ctx, from, to, draw.In, weight)
	}
}

func (drawer ImgDrawer) DrawOutArc(from draw.Pos, to draw.Pos, weight int) {
	if drawer.ctx != nil {
		draw.Arc(drawer.ctx, from, to, draw.Out, weight)
	}
}


func Png (composeNet compose.Composer) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer := ImgDrawer{draw2dimg.NewGraphicContext(img)}

	draw.Init(drawer.ctx, width, height)
	draw.Clean(drawer.ctx, width, height) // background
	composeNet(drawer)

	draw2dimg.SaveToPngFile(filename, img)
}