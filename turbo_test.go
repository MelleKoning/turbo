package turbo

import (
	"image"
	"image/color"
)

func MakeGolangRGBA() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 255, 255))

	g := uint8(0)
	for y := 0; y < 255; y++ {
		r := uint8(0)
		b := uint8(255)
		for x := 0; x < 255; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
			r += 1
			b -= 1
		}
		g += 1
	}
	return img
}
