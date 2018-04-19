package export

import (
	"git.yo2.cz/drahoslav/penego/draw"
	"github.com/llgcode/draw2d/draw2dsvg"
)

func Svg(composeNet func(draw.Drawer)) error {
	img := draw2dsvg.NewSvg()
	drawer := &ImgDrawer{draw2dsvg.NewGraphicContext(img), 0}

	width, height := store.Int("width"), store.Int("height")

	draw.Init(drawer.ctx, width, height)
	draw.Clean(drawer.ctx, width, height) // background
	composeNet(drawer)

	return draw2dsvg.SaveToSvgFile(getName("svg"), img)
}
