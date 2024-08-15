package main

import (
	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

const (
	period     = 10
	snakeWidth = 7
	segmentLen = 14
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

var snakes []*Snake

type Segment struct {
	Head firefly.Point
	Tail *Segment
}

// Render the snake's segment
func (s *Segment) Render(frame int, state State) {
	if s.Tail == nil {
		return
	}
	start := s.Head
	end := s.Tail.Head
	start.X, end.X = denormalizeX(start.X, end.X)
	start.Y, end.Y = denormalizeY(start.Y, end.Y)
	// if this is the last segment (the snake's tail), draw it shorter.
	if s.Tail.Tail == nil && state != Growing {
		end.X = start.X + (end.X-start.X)*(period-frame)/period
		end.Y = start.Y + (end.Y-start.Y)*(period-frame)/period
	}
	drawSegment(start, end)
}

type Snake struct {
	Peer firefly.Peer

	// The start point of the first full-length segment (the neck).
	Head *Segment

	// The very first point of the snake. Updated based on Dir.
	Mouth firefly.Point

	// The point the snake is looking at.
	Eye          firefly.Point
	BlinkCounter int // The timer for the snake's eye blinking.
	BlinkMaxTime int

	// The snake's movement direction in radians. Updated based on touch pad.
	Dir float32

	// Indicates if the snake is growing.
	state State
}

func NewSnake(peer firefly.Peer) *Snake {
	shift := 10 + snakeWidth + int(peer)*20
	return &Snake{
		Peer: peer,
		Head: &Segment{
			Head: firefly.Point{X: segmentLen * 2, Y: shift},
			Tail: &Segment{
				Head: firefly.Point{X: segmentLen, Y: shift},
				Tail: nil,
			},
		},
	}
}

// Update the position of all snake's segments.
func (s *Snake) Update(frame int, apple *Apple) {
	frame = frame % period
	pad, pressed := firefly.ReadPad(s.Peer)
	if pressed {
		s.setDir(pad)
	}
	if frame == 0 {
		s.shift()
	}
	s.updateMouth(frame)
	s.updateEye(apple.Pos)
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

// Make the snake look at the apple.
func (s *Snake) updateEye(apple firefly.Point) {
	// Calculate position of eye based on the where the apple is
	lookX := float32(apple.X - s.Mouth.X)
	lookY := float32(apple.Y - s.Mouth.Y)
	lookLen := tinymath.Hypot(lookX, lookY)
	dX := lookX * 3 / lookLen
	dY := lookY * 3 / lookLen

	s.Eye = firefly.Point{
		X: s.Mouth.X + int(dX),
		Y: s.Mouth.Y + int(dY),
	}

	s.BlinkCounter += int(firefly.GetRandom() % 5)
	if s.BlinkCounter > s.BlinkMaxTime {
		s.BlinkCounter = 0
		s.BlinkMaxTime = int(100 + firefly.GetRandom()%100)
	}
}

// Shift forward the position of each segment.
func (s *Snake) shift() {
	shiftX := tinymath.Cos(s.Dir) * segmentLen
	shiftY := tinymath.Sin(s.Dir) * segmentLen
	head := firefly.Point{
		X: normalizeX(s.Head.Head.X + int(shiftX)),
		Y: normalizeY(s.Head.Head.Y - int(shiftY)),
	}

	if s.state == Growing {
		s.Head = &Segment{
			Head: head,
			Tail: s.Head,
		}
		s.state = Moving
		return
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

// Update snake's mouth position based on the current frame and direction.
func (s *Snake) updateMouth(frame int) {
	neck := s.Head.Head
	headLen := float32(segmentLen) * float32(frame) / float32(period)
	shiftX := tinymath.Cos(s.Dir) * headLen
	shiftY := tinymath.Sin(s.Dir) * headLen
	x := normalizeX(neck.X + int(shiftX))
	y := normalizeY(neck.Y - int(shiftY))
	s.Mouth = firefly.Point{X: x, Y: y}
}

// Check if the snake can eat the apple.
//
// If it can, start growing the snake and move the apple.
func (s *Snake) TryEat(apple *Apple, score *Score) {
	x := apple.Pos.X - s.Mouth.X
	y := apple.Pos.Y - s.Mouth.Y
	distance := tinymath.Hypot(float32(x), float32(y))
	if distance > appleRadius+snakeWidth/2 {
		return
	}
	s.state = Eating
	apple.Move()
	score.Inc()
	// Don't place the apple inside the snake
	for s.Collides(apple.Pos) {
		apple.Move()
	}
}

// Check if the given point is within the snake's body
func (s Snake) Collides(p firefly.Point) bool {
	segment := s.Head.Tail
	for segment != nil {
		if segment.Tail != nil {
			ph := segment.Head
			pt := segment.Tail.Head
			ph.X, pt.X = denormalizeX(ph.X, pt.X)
			ph.Y, pt.Y = denormalizeY(ph.Y, pt.Y)
			bbox := NewBBox(ph, pt, snakeWidth/2)
			if bbox.Contains(p) {
				return true
			}
		}
		segment = segment.Tail
	}
	return false
}

// Render all segments and the head of the snake
func (s Snake) Render(frame int) {
	frame = frame % period
	segment := s.Head
	for segment != nil {
		segment.Render(frame, s.state)
		segment = segment.Tail
	}
	s.renderHead()
}

// Draw the zero segment of the snake: it's head.
func (s Snake) renderHead() {
	neck := s.Head.Head
	mouth := s.Mouth
	neck.X, mouth.X = denormalizeX(neck.X, mouth.X)
	neck.Y, mouth.Y = denormalizeY(neck.Y, mouth.Y)
	drawSegment(neck, mouth)
	style := firefly.Style{FillColor: firefly.ColorWhite}
	if s.Collides(mouth) {
		style.FillColor = firefly.ColorRed
	}

	firefly.DrawCircle(
		firefly.Point{
			X: mouth.X - snakeWidth/2 - 1,
			Y: mouth.Y - snakeWidth/2 - 1,
		},
		snakeWidth+2, firefly.Style{FillColor: firefly.ColorBlue},
	)
	firefly.DrawCircle(
		firefly.Point{
			X: mouth.X - snakeWidth/2,
			Y: mouth.Y - snakeWidth/2,
		},
		snakeWidth, firefly.Style{FillColor: firefly.ColorLightBlue},
	)
	firefly.DrawCircle(
		firefly.Point{
			X: s.Mouth.X - snakeWidth/2 + 1,
			Y: s.Mouth.Y - snakeWidth/2 + 1,
		},
		snakeWidth-2, style,
	)

	s.renderEye()
}

// Draw the snake's eye.
func (s Snake) renderEye() {
	firefly.DrawCircle(
		firefly.Point{
			X: s.Eye.X - snakeWidth/8,
			Y: s.Eye.Y - snakeWidth/8,
		},
		snakeWidth/4, firefly.Style{FillColor: firefly.ColorBlack},
	)

	if s.BlinkCounter < 20 {
		firefly.DrawCircle(
			firefly.Point{
				X: s.Mouth.X - snakeWidth/2 + 1,
				Y: s.Mouth.Y - snakeWidth/2 + 1,
			},
			snakeWidth-2, firefly.Style{FillColor: firefly.ColorLightBlue},
		)
	}
}

// Render the segment and ghost segments if the snake wraps around the screen edges.
func drawSegment(start, end firefly.Point) {
	drawSegmentExactlyAt(start, end)
	drawSegmentExactlyAt(
		firefly.Point{X: start.X - firefly.Width, Y: start.Y},
		firefly.Point{X: end.X - firefly.Width, Y: end.Y},
	)
	drawSegmentExactlyAt(
		firefly.Point{X: start.X, Y: start.Y - firefly.Height},
		firefly.Point{X: end.X, Y: end.Y - firefly.Height},
	)
	drawSegmentExactlyAt(
		firefly.Point{X: start.X - firefly.Width, Y: start.Y - firefly.Height},
		firefly.Point{X: end.X - firefly.Width, Y: end.Y - firefly.Height},
	)
}

// Render the segment.
func drawSegmentExactlyAt(start, end firefly.Point) {
	firefly.DrawLine(
		start, end,
		firefly.LineStyle{
			Color: firefly.ColorBlue,
			Width: snakeWidth,
		},
	)
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
