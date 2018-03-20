package export

import (
	"fmt"
	"git.yo2.cz/drahoslav/penego/draw"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
)

func Png(composeNet func(draw.Drawer)) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer := ImgDrawer{draw2dimg.NewGraphicContext(img)}

	draw.Init(drawer.ctx, width, height)
	draw.Clean(drawer.ctx, width, height) // background
	composeNet(drawer)

	err := draw2dimg.SaveToPngFile(getName("png"), img)
	if err != nil {
		fmt.Println(err)
	}
}
