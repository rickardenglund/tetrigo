package fonts

import (
	"io/ioutil"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func GetFont() font.Face {
	f, err := os.Open("assets/intuitive.ttf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fontbytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	myFont, err := truetype.Parse(fontbytes)
	if err != nil {
		panic(err)
	}

	face := truetype.NewFace(myFont, &truetype.Options{
		Size:              50,
		GlyphCacheEntries: 1,
	})

	return face
}
