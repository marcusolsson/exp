package wordcloud

import (
	"image"
	"image/color"
	"image/draw"
	"sort"
	"sync"

	"github.com/golang/freetype/truetype"
	"github.com/marcusolsson/wordcloud/quadtree"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Cloud contains the word cloud configuration.
type Cloud struct {
	Width     int
	Height    int
	Font      *truetype.Font
	Generator Generator
}

// Word contains a rendered image as well as a hierarchical bounding box used
// for intersection testing.
type Word struct {
	Text   string
	Image  image.Image
	Pos    image.Point
	Bounds *quadtree.Tree
	Weight int
}

// Intersect ...
func (w *Word) Intersect(other Word) bool {
	e1 := w.Bounds.Extents.Add(w.Pos)
	e2 := other.Bounds.Extents.Add(other.Pos)
	return e1.Overlaps(e2)
}

// Generate ...
func (c *Cloud) Generate(words map[string]int) image.Image {
	// Create canvas.
	im := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))

	// Fill canvas with background color.
	bg := &image.Uniform{color.RGBA{243, 90, 40, 255}}
	draw.Draw(im, im.Bounds(), bg, image.ZP, draw.Src)

	data := buildWords(words, c.Font)

	// Order words by weight.
	sort.Sort(ByWeight(data))

	var scene []Word
	center := image.Point{c.Width / 2, c.Height / 2}

	// Place word sprites on canvas.
	for _, w := range data {
		c.Generator.Reset()

	L:
		for {
			cb := w.Bounds.Extents

			// Generate a candidate position for the word.
			gp := c.Generator.Generate()

			w.Pos = gp.Add(center).Sub(image.Point{cb.Dx() / 2, cb.Dy() / 2})

			// Check if we can place the candidate here.
			for _, v := range scene {
				if w.Intersect(v) {
					continue L
				}
			}

			// Draw the word sprite on the canvas.
			r := image.Rectangle{w.Pos, w.Pos.Add(cb.Size())}
			draw.Draw(im, r, w.Image, image.ZP, draw.Over)

			// Add the word to the scene.
			scene = append(scene, w)

			break
		}
	}

	return im
}

// buildWords generates word sprites.
func buildWords(words map[string]int, f *truetype.Font) []Word {
	ch := make(chan Word)
	wg := &sync.WaitGroup{}
	wg.Add(len(words))

	for k, v := range words {
		go func(w string, occ int) {
			face := truetype.NewFace(f, &truetype.Options{
				Size:    float64(occ * 12),
				DPI:     72,
				Hinting: font.HintingFull,
			})

			wd := buildSprite(w, face)
			wd.Weight = occ

			ch <- wd
			wg.Done()
		}(k, v)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var data []Word
	for w := range ch {
		data = append(data, w)
	}

	return data
}

// buildSprite renders an image from a text and font face as well as computing its
// bounding box.
//
// TODO: This is a really naive implementation that really should be refactored
// into something more elegant.
func buildSprite(text string, face font.Face) Word {
	var (
		canvasWidth  = 1024
		canvasHeight = 1024
	)

	// Create temporary scratch image to determine real font size.
	scratch := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	draw.Draw(scratch, scratch.Bounds(), image.Transparent, image.ZP, draw.Src)

	d := &font.Drawer{
		Dst:  scratch,
		Src:  image.White,
		Face: face,
	}

	d.Dot = fixed.Point26_6{
		X: fixed.I(canvasWidth / 2),
		Y: fixed.I(canvasHeight / 2),
	}

	d.DrawString(text)

	// Generate trimmed image and quadtree from temporary canvas.
	tt := buildTree(scratch, color.RGBA{0, 0, 0, 0}).Trimmed()
	r := tt.Extents

	tr := image.Rect(0, 0, r.Dx(), r.Dy())
	im := image.NewRGBA(tr)
	draw.Draw(im, tr, scratch, r.Min, draw.Src)

	t := buildTree(im, color.RGBA{0, 0, 0, 0})

	return Word{
		Text:   text,
		Image:  im,
		Bounds: t,
	}
}

// buildTree is a helper function for constructing a quadtree from an image.
func buildTree(im image.Image, bg color.Color) *quadtree.Tree {
	b := im.Bounds()

	result := quadtree.New(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: b.Dx(), Y: b.Dy()},
	}, 1)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if im.At(x, y) != bg {
				result.Insert(image.Point{x, y})
			}
		}
	}

	return result
}

// drawTree is a helper function to draw a quadtree onto an image. Used for
// debugging purposes.
func drawTree(t *quadtree.Tree, m *image.RGBA, c color.Color) {
	if t == nil {
		return
	}

	b := t.Extents

	for x := b.Min.X; x < b.Max.X; x++ {
		m.Set(x, b.Min.Y, c)
	}
	for x := b.Min.X; x < b.Max.X; x++ {
		m.Set(x, b.Max.Y-1, c)
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		m.Set(b.Min.X, y, c)
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		m.Set(b.Max.X-1, y, c)
	}

	drawTree(t.NorthWest, m, c)
	drawTree(t.NorthEast, m, c)
	drawTree(t.SouthWest, m, c)
	drawTree(t.SouthEast, m, c)
}

// ByWeight sorts a slice of words by weight in descending order.
type ByWeight []Word

func (s ByWeight) Len() int {
	return len(s)
}
func (s ByWeight) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByWeight) Less(i, j int) bool {
	return s[i].Weight > s[j].Weight
}
