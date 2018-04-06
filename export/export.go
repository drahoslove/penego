package export

import (
	"path/filepath"

	"git.yo2.cz/drahoslav/penego/draw"
	"git.yo2.cz/drahoslav/penego/storage"
	"github.com/llgcode/draw2d"
)

var (
	store *storage.Storage
)

func init() {
	store = storage.Of("export")
}

type ImgDrawer struct {
	ctx   draw2d.GraphicContext
	style draw.Style
}

func (drawer *ImgDrawer) SetStyle(style draw.Style) {
	drawer.style = style
}

func (drawer ImgDrawer) DrawPlace(pos draw.Pos, n int, description string) {
	if drawer.ctx != nil {
		draw.Place(drawer.ctx, drawer.style, pos, n, description)
	}
}

func (drawer ImgDrawer) DrawTransition(pos draw.Pos, attrs, description string) {
	if drawer.ctx != nil {
		draw.Transition(drawer.ctx, drawer.style, pos, attrs, description)
	}
}

func (drawer ImgDrawer) DrawInArc(from draw.Pos, to draw.Pos, weight int) {
	if drawer.ctx != nil {
		draw.Arc(drawer.ctx, drawer.style, from, to, draw.In, weight)
	}
}

func (drawer ImgDrawer) DrawOutArc(from draw.Pos, to draw.Pos, weight int) {
	if drawer.ctx != nil {
		draw.Arc(drawer.ctx, drawer.style, from, to, draw.Out, weight)
	}
}

func getName(ext string) string {
	filename := store.Of(ext).String("filename")
	println("getName", filename, ext)
	return filename
}

func ByName(filename string, composeNet func(draw.Drawer)) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".png":
		Png(composeNet)
	case ".svg":
		Svg(composeNet)
	case ".pdf":
		Pdf(composeNet)
	}
}
