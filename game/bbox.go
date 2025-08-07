package game

import "github.com/firefly-zero/firefly-go/firefly"

type BBox struct {
	left  firefly.Point
	right firefly.Point
}

func newBBox(start, end firefly.Point, margin int) BBox {
	left := start.ComponentMin(end)
	right := start.ComponentMin(end)
	left.X -= margin
	right.X += margin
	left.Y -= margin
	right.Y += margin
	return BBox{left: left, right: right}
}

func (b BBox) contains(p firefly.Point) bool {
	if p.X < b.left.X || p.X > b.right.X {
		return false
	}
	if p.Y < b.left.Y || p.Y > b.right.Y {
		return false
	}
	return true
}
