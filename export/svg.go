package export

import (
	"fmt"

	"git.yo2.cz/drahoslav/penego/draw"
	"github.com/llgcode/draw2d/draw2dsvg"
)

func Svg(composeNet func(draw.Drawer)) {
	img := draw2dsvg.NewSvg()
	drawer := &ImgDrawer{draw2dsvg.NewGraphicContext(img), 0}

	draw.Init(drawer.ctx, width, height)
	draw.Clean(drawer.ctx, width, height) // background
	composeNet(drawer)

	err := draw2dsvg.SaveToSvgFile(getName("svg"), img)
	if err != nil {
		fmt.Println(err)
	}
}
