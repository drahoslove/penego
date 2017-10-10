package gui

// create font cache with Go font and Awesome font

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

type customFontCache map[string]*truetype.Font

func (fc customFontCache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

func (fc customFontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("Font %s is not stored in font cache.", fd.Name)
	}
	return font, nil
}

func init() {
	fontCache := customFontCache{}

	TTFs := map[string]([]byte){
		"goregular": goregular.TTF,
		"gobold":    gobold.TTF,
		"goitalic":  goitalic.TTF,
		"gomono":    gomono.TTF,
		"ico":       icoTTF,
	}

	for fontName, TTF := range TTFs {
		font, err := truetype.Parse(TTF)
		if err != nil {
			panic(err)
		}
		fontCache.Store(draw2d.FontData{Name: fontName}, font)
	}

	draw2d.SetFontCache(fontCache)
}
