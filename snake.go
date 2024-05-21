package main

import (
	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

const (
	period     = 10
	snakeWidth = 8
	segmentLen = 16
	maxDirDiff = .1
)

type State uint8

const (
	// The snake's size is stable.
	Moving State = 0

	// The snake just eat an apple and waiting to grow.
	Eating State = 1

	// The snake is growing.
	// The tail is not moving and the next shift will add a segment.
	Growing State = 2
)

var snake *Snake

type Segment struct {
	Head firefly.Point
	Tail *Segment
}

// Render the snake's segment
func (s *Segment) Render(frame int, state State) {
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
	if s.Tail.Tail == nil && state != Growing {
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

	// Indicates if the snake is growing.
	state State
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
	}
}

func (s *Snake) Update(frame int) {
	frame = frame % period
	pad, pressed := firefly.ReadPad(firefly.Player0)
	if pressed {
		s.setDir(pad)
	}
	if frame == 0 {
		s.shift()
	}
	s.updateMouth(frame)
}

// Set Dir value based on the pad input.
func (s *Snake) setDir(pad firefly.Pad) {
	dirDiff := pad.Azimuth().Radians() - s.Dir
	if tinymath.IsNaN(dirDiff) {
		return
	}

	// If the turn is more than 180 degrees, we're rotating in a wrong direction.
	// Switch the direction.
	if dirDiff > tinymath.Pi {
		dirDiff = -maxDirDiff
	}
	if dirDiff < -tinymath.Pi {
		dirDiff = maxDirDiff
	}

	// Smoothen the turn.
	if dirDiff > maxDirDiff {
		s.Dir += maxDirDiff
	} else if dirDiff < -maxDirDiff {
		s.Dir -= maxDirDiff
	} else {
		s.Dir += dirDiff
	}

	// Ensure that the direction is always on the 0-360 degrees range.
	if s.Dir < 0 {
		s.Dir = tinymath.Tau - s.Dir
	}
	if s.Dir > tinymath.Tau {
		s.Dir = s.Dir - tinymath.Tau
	}
}

// Shift forward the position of each segment.
func (s *Snake) shift() {
	shiftX := tinymath.Cos(s.Dir) * segmentLen
	shiftY := tinymath.Sin(s.Dir) * segmentLen
	head := firefly.Point{
		X: s.Head.Head.X + int(shiftX),
		Y: s.Head.Head.Y - int(shiftY),
	}

	if s.state == Growing {
		s.Head = &Segment{
			Head: head,
			Tail: s.Head,
		}
		s.state = Moving
	}
	if s.state == Eating {
		s.state = Growing
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
	headLen := float32(segmentLen) * float32(frame) / float32(period)
	shiftX := tinymath.Cos(s.Dir) * headLen
	shiftY := tinymath.Sin(s.Dir) * headLen
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
	distance := tinymath.Hypot(float32(x), float32(y))
	if distance > appleRadius+snakeWidth/2 {
		return
	}
	s.state = Eating
	a.Move()
}

func (s *Snake) Render(frame int) {
	frame = frame % period
	segment := s.Head
	for segment != nil {
		segment.Render(frame, s.state)
		segment = segment.Tail
	}
	s.renderHead()
}

// Draw the zero segment of the snake: it's head.
func (s *Snake) renderHead() {
	neck := s.Head.Head
	lineStyle := firefly.LineStyle{
		Color: firefly.ColorBlue,
		Width: snakeWidth,
	}
	firefly.DrawLine(neck, s.Mouth, lineStyle)

	style := firefly.Style{FillColor: firefly.ColorBlue}
	firefly.DrawCircle(
		firefly.Point{
			X: neck.X - snakeWidth/2,
			Y: neck.Y - snakeWidth/2,
		},
		snakeWidth, style,
	)
	firefly.DrawCircle(
		firefly.Point{
			X: s.Mouth.X - snakeWidth/2,
			Y: s.Mouth.Y - snakeWidth/2,
		},
		snakeWidth, style,
	)
}
