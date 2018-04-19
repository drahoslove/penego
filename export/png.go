package export

import (
	"image"

	"git.yo2.cz/drahoslav/penego/draw"
	"github.com/llgcode/draw2d/draw2dimg"
)

func Png(composeNet func(draw.Drawer)) error {
	width, height := store.Int("width"), store.Int("height")

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer := &ImgDrawer{draw2dimg.NewGraphicContext(img), 0}

	draw.Init(drawer.ctx, width, height)
	draw.Clean(drawer.ctx, width, height) // background
	composeNet(drawer)

	return draw2dimg.SaveToPngFile(getName("png"), img)
}
