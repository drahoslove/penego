package export

import (
	"fmt"
	"git.yo2.cz/drahoslav/penego/draw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dpdf"
)

func Pdf(composeNet func(draw.Drawer)) {
	img := draw2dpdf.NewPdf("L", "mm", "A4")
	drawer := ImgDrawer{draw2dpdf.NewGraphicContext(img)}

	draw2d.SetFontNamer(func(fd draw2d.FontData) string {
		return fd.Name
	})
	draw.Init(drawer.ctx, width, height)
	drawer.ctx.SetFontData(draw2d.FontData{Name: "courier"}) // TODO use gomono
	draw.Clean(drawer.ctx, width, height)                    // background
	composeNet(drawer)

	err := draw2dpdf.SaveToPdfFile(getName("pdf"), img)
	if err != nil {
		fmt.Println(err)
	}
}
