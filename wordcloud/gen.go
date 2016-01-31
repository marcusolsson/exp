package wordcloud

import (
	"image"
	"math"
)

// Generator ...
type Generator interface {
	// Reset ...
	Reset()

	// Generate ...
	Generate() image.Point
}

// SpiralGenerator ...
type SpiralGenerator struct {
	radius float64
	theta  float64
}

// Reset ...
func (g *SpiralGenerator) Reset() {
	g.radius = 0.0
	g.theta = 0.0
}

// Generate ...
func (g *SpiralGenerator) Generate() image.Point {
	x, y := int(g.radius*math.Sin(g.theta)), int(g.radius*math.Cos(g.theta))

	g.radius += 0.05
	g.theta += math.Pi / 100.0

	return image.Point{x, y}
}

// NewSpiralGenerator ...
func NewSpiralGenerator() Generator {
	return &SpiralGenerator{}
}
