package utils

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"

	"github.com/RobCherry/vibrant"
	"github.com/nfnt/resize"
)

func getImagefromURL(url string) (error, *image.Image) {
	resp, err := http.Get(url)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	m, _, err := image.Decode(resp.Body)
	if err != nil {
		return err, nil
	}
	return nil, &m
}

func getPalettefromURL(url string) (error, *vibrant.Palette) {
	err, image := getImagefromURL(url)
	if err != nil {
		return err, nil
	}
	g := (*image).Bounds()
	// Get height and width
	height := float64(g.Dy())
	width := float64(g.Dx())
	scaleRatio := 100 / float64(math.Min(width, height))
	m := resize.Resize(uint(width*scaleRatio), 0, *image, resize.Lanczos3)
	paletteBuilder := vibrant.NewPaletteBuilder(m).MaximumColorCount(16)
	return nil, paletteBuilder.Generate()
}
