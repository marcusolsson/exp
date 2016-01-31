package quadtree

import "image"

// Contains ...
func Contains(bb image.Rectangle, p image.Point) bool {
	if bb.Min.X <= p.X && p.X < bb.Max.X {
		if bb.Min.Y <= p.Y && p.Y < bb.Max.Y {
			return true
		}
	}
	return false
}

// Tree ...
type Tree struct {
	Extents image.Rectangle

	points   map[image.Point]struct{}
	capacity int

	NorthWest *Tree
	NorthEast *Tree
	SouthWest *Tree
	SouthEast *Tree
}

// New returns new instance of Tree.
func New(r image.Rectangle, c int) *Tree {
	return &Tree{
		points:   make(map[image.Point]struct{}),
		capacity: c,
		Extents:  r,
	}
}

// Insert ...
func (t *Tree) Insert(p image.Point) bool {
	if !Contains(t.Extents, p) {
		return false
	}

	if len(t.points) < t.capacity {
		t.points[p] = struct{}{}
		return true
	}

	if t.NorthWest == nil {
		t.Subdivide()

		for p := range t.points {
			t.Insert(p)
		}
		t.points = make(map[image.Point]struct{})
	}

	if t.NorthWest.Insert(p) {
		return true
	}
	if t.NorthEast.Insert(p) {
		return true
	}
	if t.SouthWest.Insert(p) {
		return true
	}
	if t.SouthEast.Insert(p) {
		return true
	}

	return false
}

// Subdivide ...
func (t *Tree) Subdivide() {
	b := t.Extents
	c := b.Min.Add(image.Point{b.Dx() / 2, b.Dy() / 2})

	var (
		southWest = image.Point{b.Min.X, b.Min.Y}
		west      = image.Point{b.Min.X, c.Y}
		east      = image.Point{b.Max.X, c.Y}
		northEast = image.Point{b.Max.X, b.Max.Y}
		south     = image.Point{c.X, b.Min.Y}
		north     = image.Point{c.X, b.Max.Y}
	)

	t.NorthWest = New(image.Rectangle{Min: west, Max: north}, t.capacity)
	t.NorthEast = New(image.Rectangle{Min: c, Max: northEast}, t.capacity)
	t.SouthWest = New(image.Rectangle{Min: southWest, Max: c}, t.capacity)
	t.SouthEast = New(image.Rectangle{Min: south, Max: east}, t.capacity)
}

// Trimmed ...
func (t *Tree) Trimmed() *Tree {
	var leafs []image.Rectangle
	t.Visit(func(other *Tree) {
		if len(other.points) == 0 {
			return
		}
		leafs = append(leafs, image.Rectangle{
			Min: other.Extents.Min,
			Max: other.Extents.Max,
		})
	})

	var r *image.Rectangle
	for _, b := range leafs {
		if r == nil {
			r = new(image.Rectangle)
			*r = b
			continue
		}
		*r = r.Union(b)
	}

	tr := New(*r, t.capacity)

	t.Visit(func(other *Tree) {
		for p := range other.points {
			tr.Insert(p)
		}
	})

	return tr
}

// Visit ...
func (t *Tree) Visit(fn func(*Tree)) {
	if t.NorthWest == nil {
		fn(t)
		return
	}
	t.NorthWest.Visit(fn)
	t.NorthEast.Visit(fn)
	t.SouthWest.Visit(fn)
	t.SouthEast.Visit(fn)
}
