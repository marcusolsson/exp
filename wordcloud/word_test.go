package wordcloud

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var defaultFont = "/usr/share/fonts/truetype/roboto/Roboto-Regular.ttf"

func BenchmarkBuildWords(b *testing.B) {
	f, err := readFont(defaultFont)
	if err != nil {
		b.FailNow()
	}

	for i := 0; i < b.N; i++ {
		buildWords(map[string]int{
			"test1": 1,
			"test2": 2,
			"test3": 3,
			"test4": 4,
			"test5": 5,
			"test6": 6,
		}, f)
	}
}

func BenchmarkBuildSprite(b *testing.B) {
	f, err := readFont(defaultFont)
	if err != nil {
		b.FailNow()
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:    12.0,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	for i := 0; i < b.N; i++ {
		buildSprite("test", face)
	}
}

func BenchmarkBuildTree(b *testing.B) {
	im := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(im, im.Bounds(), image.White, image.ZP, draw.Src)

	bb := im.Bounds()

	for y := bb.Min.Y; y < bb.Max.Y; y++ {
		for x := bb.Min.X; x < bb.Max.X; x++ {
			i := x + bb.Dx()/2
			j := y + bb.Dy()/2

			r := math.Sqrt(float64(i*i + j*j))

			if r < 200 {
				im.Set(x, y, color.Black)
			}
		}
	}

	for i := 0; i < b.N; i++ {
		buildTree(im, color.White)
	}
}

func BenchmarkGenerateCloud(b *testing.B) {
	f, err := readFont(defaultFont)
	if err != nil {
		b.FailNow()
	}

	c := Cloud{
		Width:     1024,
		Height:    1024,
		Font:      f,
		Generator: NewSpiralGenerator(),
	}

	for i := 0; i < b.N; i++ {
		c.Generate(map[string]int{
			"test1": 1,
			"test2": 2,
			"test3": 3,
			"test4": 4,
			"test5": 5,
			"test6": 6,
		})
	}
}

func readFont(s string) (*truetype.Font, error) {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}

	return f, nil
}
