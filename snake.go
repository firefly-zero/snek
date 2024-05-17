package main

import "github.com/firefly-zero/firefly-go/firefly"

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
			X: end.X - snakeWidth/2 + 1,
			Y: end.Y - snakeWidth/2 + 1,
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

func (s *Snake) Render(frame int) {
	frame = frame % period
	segment := s.Head
	for segment != nil {
		segment.Render(frame)
		segment = segment.Tail
	}
}
