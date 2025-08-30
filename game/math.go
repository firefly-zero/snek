package game

import (
	"unsafe"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Line struct {
	h firefly.Point
	t firefly.Point
}

// Check if two segments intersect.
//
// Source: https://bryceboe.com/2006/10/23/line-segment-intersection-algorithm/
func intersect(a, b Line) bool {
	return ccw(a.h, b.h, b.t) != ccw(a.t, b.h, b.t) && ccw(a.h, a.t, b.h) != ccw(a.h, a.t, b.t)
}

// Check if the two points are in the counter-closckwise order.
func ccw(a, b, c firefly.Point) bool {
	return (c.Y-a.Y)*(b.X-a.X) > (b.Y-a.Y)*(c.X-a.X)
}

// If x points outside the screen, shift it so that it's back on the screen.
func normalizeX(x int) int {
	if x >= firefly.Width {
		x -= firefly.Width
	} else if x < 0 {
		x += firefly.Width
	}
	return x
}

// If y points outside the screen, shift it so that it's back on the screen.
func normalizeY(y int) int {
	if y >= firefly.Height {
		y = y - firefly.Height
	} else if y < 0 {
		y += firefly.Height
	}
	return y
}

// If the dots are on the opposite sides of the screen,
// put the left one on the right outside the screen.
func denormalizeX(start, end int) (int, int) {
	if start-end > 30 {
		end += firefly.Width
	} else if end-start > 30 {
		start += firefly.Width
	}
	return start, end
}

// If the dots are on the opposite sides of the screen,
// put the upper one on the bottom outside the screen.
func denormalizeY(start, end int) (int, int) {
	if start-end > 30 {
		end += firefly.Height
	} else if end-start > 30 {
		start += firefly.Height
	}
	return start, end
}

func formatInt(i int16) string {
	buf := []byte{'0' + byte(i/10), '0' + byte(i%10)}
	return unsafe.String(&buf[0], 2)
}
