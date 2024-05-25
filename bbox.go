package main

import "github.com/firefly-zero/firefly-go/firefly"

type BBox struct {
	left  firefly.Point
	right firefly.Point
}

func NewBBox(start, end firefly.Point, margin int) BBox {
	if end.X < start.X {
		start.X, end.X = end.X, start.X
	}
	if end.Y < start.Y {
		start.Y, end.Y = end.Y, start.Y
	}
	start.X -= margin
	end.X += margin
	start.Y -= margin
	end.Y += margin
	return BBox{left: start, right: end}
}

func (b BBox) Contains(p firefly.Point) bool {
	if p.X < b.left.X || p.X > b.right.X {
		return false
	}
	if p.Y < b.left.Y || p.Y > b.right.Y {
		return false
	}
	return true
}
