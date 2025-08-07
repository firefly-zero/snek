package game

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
	moving State = 0

	// The snake just eat an apple and waiting to grow.
	eating State = 1

	// The snake is growing.
	// The tail is not moving and the next shift will add a segment.
	growing State = 2
)

var snakes []*Snake

type Segment struct {
	head firefly.Point
	tail *Segment
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

type Snake struct {
	peer firefly.Peer

	// The start point of the first full-length segment (the neck).
	head *Segment

	// The very first point of the snake. Updated based on Dir.
	mouth firefly.Point

	// The point the snake is looking at.
	eye firefly.Point
	// The timer for the snake's eye blinking.
	blinkCounter int
	blinkMaxTime int

	// The snake's movement direction in radians. Updated based on touch pad.
	dir float32

	// Indicates if the snake is growing.
	state State
}

func newSnake(peer firefly.Peer) *Snake {
	shift := 10 + snakeWidth + int(peer)*20
	return &Snake{
		peer: peer,
		head: &Segment{
			head: firefly.Point{X: segmentLen * 2, Y: shift},
			tail: &Segment{
				head: firefly.Point{X: segmentLen, Y: shift},
				tail: nil,
			},
		},
	}
}

// update the position of all snake's segments.
func (s *Snake) update(frame int, apple *Apple) {
	frame = frame % period
	pad, pressed := firefly.ReadPad(s.peer)
	if pressed {
		s.setDir(pad)
	}
	if frame == 0 {
		s.shift()
	}
	s.updateMouth(frame)
	s.updateEye(apple.pos)
}

// Set Dir value based on the pad input.
func (s *Snake) setDir(pad firefly.Pad) {
	dirDiff := pad.Azimuth().Radians() - s.dir
	if tinymath.IsNaN(dirDiff) {
		return
	}

	// If the turn is more than 180 degrees, we're rotating in a wrong direction.
	// Switch the direction.
	if dirDiff > tinymath.Pi {
		dirDiff = -maxDirDiff
	} else if dirDiff < -tinymath.Pi {
		dirDiff = maxDirDiff
	}

	// Smoothen the turn.
	if dirDiff > maxDirDiff {
		s.dir += maxDirDiff
	} else if dirDiff < -maxDirDiff {
		s.dir -= maxDirDiff
	} else {
		s.dir += dirDiff
	}

	// Ensure that the direction is always on the 0-360 degrees range.
	if s.dir < 0 {
		s.dir = s.dir + tinymath.Tau
	}
	if s.dir > tinymath.Tau {
		s.dir = s.dir - tinymath.Tau
	}
}

// Make the snake look at the apple.
func (s *Snake) updateEye(apple firefly.Point) {
	// Calculate position of eye based on the where the apple is
	lookX := float32(apple.X - s.mouth.X)
	lookY := float32(apple.Y - s.mouth.Y)
	lookLen := tinymath.Hypot(lookX, lookY)
	dX := lookX * 3 / lookLen
	dY := lookY * 3 / lookLen

	s.eye = firefly.Point{
		X: s.mouth.X + int(dX),
		Y: s.mouth.Y + int(dY),
	}

	s.blinkCounter += int(firefly.GetRandom() % 5)
	if s.blinkCounter > s.blinkMaxTime {
		s.blinkCounter = 0
		s.blinkMaxTime = int(100 + firefly.GetRandom()%100)
	}
}

// Shift forward the position of each segment.
func (s *Snake) shift() {
	shiftX := tinymath.Cos(s.dir) * segmentLen
	shiftY := tinymath.Sin(s.dir) * segmentLen
	head := firefly.Point{
		X: normalizeX(s.head.head.X + int(shiftX)),
		Y: normalizeY(s.head.head.Y - int(shiftY)),
	}

	if s.state == growing {
		s.head = &Segment{
			head: head,
			tail: s.head,
		}
		s.state = moving
		return
	}
	if s.state == eating {
		s.state = growing
	}

	segment := s.head
	for segment != nil {
		oldHead := segment.head
		segment.head = head
		head = oldHead
		segment = segment.tail
	}
}

// Update snake's mouth position based on the current frame and direction.
func (s *Snake) updateMouth(frame int) {
	neck := s.head.head
	headLen := float32(segmentLen) * float32(frame) / float32(period)
	shiftX := tinymath.Cos(s.dir) * headLen
	shiftY := tinymath.Sin(s.dir) * headLen
	x := normalizeX(neck.X + int(shiftX))
	y := normalizeY(neck.Y - int(shiftY))
	s.mouth = firefly.Point{X: x, Y: y}
}

// Check if the snake can eat the apple.
//
// If it can, start growing the snake and move the apple.
func (s *Snake) tryEat(apple *Apple, score *Score) {
	x := apple.pos.X - s.mouth.X
	y := apple.pos.Y - s.mouth.Y
	distance := tinymath.Hypot(float32(x), float32(y))
	if distance > appleRadius+snakeWidth/2 {
		return
	}
	s.state = eating
	apple.move()
	score.inc()
	// Don't place the apple inside the snake
	for s.collides(apple.pos) {
		apple.move()
	}
}

// Check if the given point is within the snake's body
func (s Snake) collides(p firefly.Point) bool {
	segment := s.head.tail
	for segment != nil {
		if segment.tail != nil {
			ph := segment.head
			pt := segment.tail.head
			ph.X, pt.X = denormalizeX(ph.X, pt.X)
			ph.Y, pt.Y = denormalizeY(ph.Y, pt.Y)
			bbox := newBBox(ph, pt, snakeWidth/2)
			if bbox.contains(p) {
				return true
			}
		}
		segment = segment.tail
	}
	return false
}

// render all segments and the head of the snake
func (s Snake) render(frame int) {
	frame = frame % period
	segment := s.head
	for segment != nil {
		segment.render(frame, s.state)
		segment = segment.tail
	}
	s.renderHead()
}

// Draw the zero segment of the snake: it's head.
func (s Snake) renderHead() {
	neck := s.head.head
	mouth := s.mouth
	neck.X, mouth.X = denormalizeX(neck.X, mouth.X)
	neck.Y, mouth.Y = denormalizeY(neck.Y, mouth.Y)
	drawSegment(neck, mouth)
	style := firefly.Style{FillColor: firefly.ColorWhite}
	if s.collides(mouth) {
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
			X: s.mouth.X - snakeWidth/2 + 1,
			Y: s.mouth.Y - snakeWidth/2 + 1,
		},
		snakeWidth-2, style,
	)

	s.renderEye()
}

// Draw the snake's eye.
func (s Snake) renderEye() {
	firefly.DrawCircle(
		firefly.Point{
			X: s.eye.X - snakeWidth/8,
			Y: s.eye.Y - snakeWidth/8,
		},
		snakeWidth/4, firefly.Style{FillColor: firefly.ColorBlack},
	)

	if s.blinkCounter < 20 {
		firefly.DrawCircle(
			firefly.Point{
				X: s.mouth.X - snakeWidth/2 + 1,
				Y: s.mouth.Y - snakeWidth/2 + 1,
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
