package game

import "github.com/firefly-zero/firefly-go/firefly"

type Segment struct {
	head firefly.Point
	tail *Segment
}

func (s *Segment) line() Line {
	ph := s.head
	pt := s.tail.head
	ph.X, pt.X = denormalizeX(ph.X, pt.X)
	ph.Y, pt.Y = denormalizeY(ph.Y, pt.Y)
	return Line{ph, pt}
}

// render the snake's segment
func (s *Segment) render(frame int, state State) {
	if s.tail == nil {
		return
	}
	start := s.head
	end := s.tail.head
	start.X, end.X = denormalizeX(start.X, end.X)
	start.Y, end.Y = denormalizeY(start.Y, end.Y)
	// if this is the last segment (the snake's tail), draw it shorter.
	if s.tail.tail == nil && state != growing {
		end.X = start.X + (end.X-start.X)*(period-frame)/period
		end.Y = start.Y + (end.Y-start.Y)*(period-frame)/period
	}
	drawSegment(start, end)
}
