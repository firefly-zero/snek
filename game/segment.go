package game

import "github.com/firefly-zero/firefly-go/firefly"

type Segment struct {
	head firefly.Point
	tail *Segment
	hurt bool
}

func (s *Segment) line() Line {
	ph := s.head
	pt := s.tail.head
	ph.X, pt.X = denormalizeX(ph.X, pt.X)
	ph.Y, pt.Y = denormalizeY(ph.Y, pt.Y)
	return Line{ph, pt}
}

// render the snake's segment
func (s *Segment) render(frame int, state State, me bool) {
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
	c := firefly.ColorRed
	if !s.hurt {
		switch ((s.head.X + s.head.Y) / 12) % 3 {
		case 0:
			if me {
				c = firefly.ColorDarkBlue
			} else {
				c = firefly.ColorDarkGray
			}
		case 1:
			if me {
				c = firefly.ColorLightBlue
			} else {
				c = firefly.ColorLightGray
			}
		case 2:
			if me {
				c = firefly.ColorBlue
			} else {
				c = firefly.ColorGray
			}
		}
	}
	drawSegment(start, end, c)
}
