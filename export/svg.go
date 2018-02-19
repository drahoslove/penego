package export

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dsvg"
	"git.yo2.cz/drahoslav/penego/draw"
)

func Svg (composeNet func(draw.Drawer)) {
	img := draw2dsvg.NewSvg()
	drawer := ImgDrawer{draw2dsvg.NewGraphicContext(img)}

	draw.Init(drawer.ctx, width, height)
	draw.Clean(drawer.ctx, width, height) // background
	composeNet(drawer)

	err := draw2dsvg.SaveToSvgFile(getName("svg"), img)
	if err != nil {
		fmt.Println(err)
	}
}
