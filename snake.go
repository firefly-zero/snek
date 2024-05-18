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
	Head *Segment
	Dir  float32
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
	if frame == period-1 {
		s.Shift()
	}
}

func (s *Snake) Shift() {
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

func (s *Snake) Render(frame int) {
	frame = frame % period
	segment := s.Head
	for segment != nil {
		segment.Render(frame)
		segment = segment.Tail
	}
	s.renderHead(frame)
}

func (s *Snake) renderHead(frame int) {
	neck := s.Head.Head
	headLen := float64(segmentLen) / float64(period) * float64(frame)
	shiftX := math.Cos(float64(s.Dir)) * float64(headLen)
	shiftY := math.Sin(float64(s.Dir)) * float64(headLen)
	mouth := firefly.Point{
		X: neck.X + int(shiftX),
		Y: neck.Y - int(shiftY),
	}
	style := firefly.LineStyle{
		Color: firefly.ColorBlue,
		Width: snakeWidth,
	}
	firefly.DrawLine(neck, mouth, style)
}
