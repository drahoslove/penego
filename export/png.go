package export

import (
	"image"
	"github.com/llgcode/draw2d/draw2dimg"
	"git.yo2.cz/drahoslav/penego/draw"
)

func Png (composeNet func(draw.Drawer)) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer := ImgDrawer{draw2dimg.NewGraphicContext(img)}

	draw.Init(drawer.ctx, width, height)
	draw.Clean(drawer.ctx, width, height) // background
	composeNet(drawer)

	draw2dimg.SaveToPngFile(getName("png"), img)
}
