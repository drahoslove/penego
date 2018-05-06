package export

import (
	"git.yo2.cz/drahoslav/penego/draw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dpdf"
)

func Pdf(composeNet func(draw.Drawer)) error {
	draw2d.SetFontFolder("fonts")
	draw2d.SetFontNamer(func(fd draw2d.FontData) string {
		return fd.Name
	})

	width, height := store.Int("width"), store.Int("height")

	orientation := ""
	if width > height {
		orientation = "L"
	} else {
		orientation = "P"
	}

	paper := ""
	switch width*height {
	case 1190*841:
		paper = "a3"
	case 842*595:
		paper = "a4"
	case 595 *420:
		paper = "a5"
	}

	img := draw2dpdf.NewPdf(orientation, "pt", paper)
	drawer := &ImgDrawer{draw2dpdf.NewGraphicContext(img), 0}

	drawer.ctx.Save() // required for pdf backend to call save before translate
	draw.Init(drawer.ctx, width, height)
	// drawer.ctx.SetFontData(draw2d.FontData{Name: "courier"}) // TODO use gomono
	draw.Clean(drawer.ctx, width, height) // background
	composeNet(drawer)
	drawer.ctx.Restore()
	return draw2dpdf.SaveToPdfFile(getName("pdf"), img)
}
