// Package export defines ImageDrawer - implementation of draw.Drawer interface
// which support drawing to various image formats
// So far following are supported: PNG, SVG, PDF
package export

import (
	"fmt"
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
func (drawer ImgDrawer) DrawInhibitorArc(from, to draw.Pos) {
	if drawer.ctx != nil {
		draw.InhibitorArc(drawer.ctx, drawer.style, from, to)
	}
}

func getName(ext string) string {
	filename := store.Of(ext).String("filename")
	return filename
}

func setName(filename string) string {
	ext := filepath.Ext(filename)[1:]
	store.Of(ext).Set("filename", filename)
	return ext
}

func ByName(filename string, composeNet func(draw.Drawer)) error {
	ext := setName(filename)
	switch ext {
	case "png":
		return Png(composeNet)
	case "svg":
		return Svg(composeNet)
	case "pdf":
		return Pdf(composeNet)
	default:
		return fmt.Errorf("Unknown export format %s", ext)
	}
}
