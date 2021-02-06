package fonts

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io/ioutil"
	"os"
)

func GetFont() font.Face {
	f, err := os.Open("assets/intuitive.ttf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//fontbytes := []byte{}
	//_, err = f.Read(fontbytes)
	fontbytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(fontbytes)
	if err != nil {
		panic(err)
	}

	face := truetype.NewFace(font, &truetype.Options{
		Size:              50,
		GlyphCacheEntries: 1,
	})
	return face
}
