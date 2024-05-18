package main

import (
	"math"

	"github.com/firefly-zero/firefly-go/firefly"
)

const (
	period     = 10
	snakeWidth = 8
	segmentLen = 16
)

var snake *Snake

type Segment struct {
	Head firefly.Point
	Tail *Segment
}

// Render the snake's segment
func (s *Segment) Render(frame int) {
	style := firefly.LineStyle{
		Color: firefly.ColorBlue,
		Width: snakeWidth,
	}
	if s.Tail == nil {
		return
	}
	start := s.Head
	end := s.Tail.Head
	// if this is the last segment (the snake's tail), draw it shorter.
	if s.Tail.Tail == nil {
		end.X = start.X + ((end.X - start.X) * (period - frame) / period)
		end.Y = start.Y + ((end.Y - start.Y) * (period - frame) / period)
	}
	firefly.DrawLine(start, end, style)
	firefly.DrawCircle(
		firefly.Point{
			X: end.X - snakeWidth/2,
			Y: end.Y - snakeWidth/2,
		},
		snakeWidth,
		firefly.Style{
			FillColor: firefly.ColorBlue,
		},
	)
}

type Snake struct {
	// The start point of the first full-length segment (the neck).
	Head *Segment
	// The very first point of the snake. Updated based on Dir.
	Mouth firefly.Point
	// The snake's movement direction in radians. Updated based on touch pad.
	Dir float32
	// Indicates that the snake eat an apple and is currently growing.
	growing bool
}

func NewSnake() *Snake {
	return &Snake{
		Head: &Segment{
			Head: firefly.Point{X: segmentLen * 2, Y: 10 + snakeWidth},
			Tail: &Segment{
				Head: firefly.Point{X: segmentLen, Y: 10 + snakeWidth},
				Tail: nil,
			},
		},
		Dir: 0,
	}
}

func (s *Snake) Update(frame int) {
	frame = frame % period
	pad, pressed := firefly.ReadPad(firefly.Player0)
	if pressed {
		s.Dir = pad.Azimuth().Radians()
	}
	if frame == 0 {
		s.shift()
	}
	s.updateMouth(frame)
}

// Shift forward the position of each segment.
func (s *Snake) shift() {
	shiftX := math.Cos(float64(s.Dir)) * float64(segmentLen)
	shiftY := math.Sin(float64(s.Dir)) * float64(segmentLen)
	head := firefly.Point{
		X: s.Head.Head.X + int(shiftX),
		Y: s.Head.Head.Y - int(shiftY),
	}
	segment := s.Head
	for segment != nil {
		oldHead := segment.Head
		segment.Head = head
		head = oldHead
		segment = segment.Tail
	}
}

func (s *Snake) updateMouth(frame int) {
	neck := s.Head.Head
	headLen := float64(segmentLen) * float64(frame) / float64(period)
	shiftX := math.Cos(float64(s.Dir)) * headLen
	shiftY := math.Sin(float64(s.Dir)) * headLen
	s.Mouth = firefly.Point{
		X: neck.X + int(shiftX),
		Y: neck.Y - int(shiftY),
	}
}

// Check if the snake can eat the apple.
//
// If it can, start growing the snake and move the apple.
func (s *Snake) TryEat(a *Apple) {
	x := a.Pos.X - s.Mouth.X
	y := a.Pos.Y - s.Mouth.Y
	distance := math.Sqrt(float64(x*x + y*y))
	if distance > appleRadius {
		return
	}
	s.growing = true
	a.Move()
}

func (s *Snake) Render(frame int) {
	frame = frame % period
	segment := s.Head
	for segment != nil {
		segment.Render(frame)
		segment = segment.Tail
	}
	s.renderHead()
}

// Draw the zero segment of the snake: it's head.
func (s *Snake) renderHead() {
	style := firefly.LineStyle{
		Color: firefly.ColorBlue,
		Width: snakeWidth,
	}
	firefly.DrawLine(s.Head.Head, s.Mouth, style)
}
