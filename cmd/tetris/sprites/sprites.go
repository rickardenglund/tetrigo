package sprites

import (
	"image"
	"os"

	"github.com/faiface/pixel"
)

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func GetBlockSprites(sheet pixel.Picture) []*pixel.Sprite {
	const spriteWidth = 64
	sprites := make([]*pixel.Sprite, 0, int(sheet.Bounds().W()/spriteWidth))
	for x := 0.0; x < sheet.Bounds().W(); x += spriteWidth {
		sprite := pixel.NewSprite(sheet, pixel.R(x, 0.0, x+spriteWidth, sheet.Bounds().H()))
		sprites = append(sprites, sprite)
	}

	return sprites
}
