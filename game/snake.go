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

type Snake struct {
	peer firefly.Peer

	score *Score

	// The start point of the first full-length segment (the neck).
	head *Segment

	// The very first point of the snake. Updated based on Dir.
	mouth firefly.Point

	eye Eye

	// The snake's movement direction in radians. Updated based on touch pad.
	dir float32

	// Indicates if the snake is growing.
	state State

	// While not zero, an "you" message will be shown above the snake's head.
	youTTL uint8
}

func newSnake(peer firefly.Peer) *Snake {
	shift := 10 + snakeWidth + int(peer)*20
	var youTTL uint8
	if peer == me {
		youTTL = 180
	}
	return &Snake{
		peer:   peer,
		score:  newScore(peer),
		youTTL: youTTL,
		eye:    Eye{peer: peer},
		head: &Segment{
			head: firefly.P(segmentLen*2, shift),
			tail: &Segment{
				head: firefly.P(segmentLen, shift),
				tail: nil,
			},
		},
	}
}

// update the position of all snake's segments.
func (s *Snake) update() {
	frame = frame % period
	if s.youTTL > 0 {
		s.youTTL--
	}
	pad, pressed := firefly.ReadPad(s.peer)
	if pressed {
		s.setDir(pad)
	}
	if frame == 0 {
		s.shift()
	}
	s.updateMouth(frame)
	s.eye.update(s.mouth)
	s.score.update()
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
		s.dir += tinymath.Tau
	}
	if s.dir > tinymath.Tau {
		s.dir -= tinymath.Tau
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
	s.mouth = firefly.P(x, y)
}

// Check if the snake can eat the apple.
//
// If it can, start growing the snake and move the apple.
func (s *Snake) tryEat() {
	const minDist = (appleRadius+snakeWidth)/2 + 3
	const minDist2 = minDist * minDist
	x := float32(apple.pos.X - s.mouth.X)
	y := float32(apple.pos.Y - s.mouth.Y)
	dist2 := x*x + y*y
	if dist2 > minDist2 {
		return
	}
	s.state = eating
	apple.move()
	s.score.inc()
}

// Check if the given apple position is within the snake's body.
func (s *Snake) appleCollides(p firefly.Point) bool {
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

// Check if this snake bites the given snake.
//
// Bites is detected based on if the first segment of this snake
// intersects any of the segments of the other snake.
func (s *Snake) bites(me bool, other *Snake) bool {
	neck := &Segment{head: s.mouth, tail: s.head}
	neckLine := neck.line()
	segment := other.head.tail
	if segment != nil && me {
		segment = segment.tail
	}
	for segment != nil {
		if segment.tail != nil {
			segment.hurt = false
			if intersect(segment.line(), neckLine) {
				segment.hurt = true
				return true
			}
		}
		segment = segment.tail
	}
	return false
}

// render all segments and the head of the snake
func (s *Snake) render() {
	frame = frame % period
	segment := s.head
	for segment != nil {
		segment.render(frame, s.state)
		segment = segment.tail
	}
	s.renderNeck()
	s.eye.render(s.mouth)
	if s.youTTL != 0 {
		s.renderYou()
	} else if s.score.ttl != 0 {
		s.renderScore()
	}
}

// Draw the zero segment of the snake: it's neck.
func (s *Snake) renderNeck() {
	neck := s.head.head
	mouth := s.mouth
	neck.X, mouth.X = denormalizeX(neck.X, mouth.X)
	neck.Y, mouth.Y = denormalizeY(neck.Y, mouth.Y)
	drawSegment(neck, mouth, firefly.ColorBlue)
}

// Render a "you" message above the snake's head.
func (s *Snake) renderYou() {
	p := firefly.P(s.mouth.X-5, s.mouth.Y-6)
	font.Draw("you", p, firefly.ColorRed)
}

func (s *Snake) renderScore() {
	font.Draw(
		formatInt(s.score.val),
		firefly.P(s.mouth.X-5, s.mouth.Y-6),
		s.score.color,
	)
}

// Render the segment and ghost segments if the snake wraps around the screen edges.
func drawSegment(start, end firefly.Point, c firefly.Color) {
	drawSegmentExactlyAt(start, end, c)
	drawSegmentExactlyAt(
		firefly.P(start.X-firefly.Width, start.Y),
		firefly.P(end.X-firefly.Width, end.Y),
		c,
	)
	drawSegmentExactlyAt(
		firefly.P(start.X, start.Y-firefly.Height),
		firefly.P(end.X, end.Y-firefly.Height),
		c,
	)
	drawSegmentExactlyAt(
		firefly.P(start.X-firefly.Width, start.Y-firefly.Height),
		firefly.P(end.X-firefly.Width, end.Y-firefly.Height),
		c,
	)
}

// Render the segment.
func drawSegmentExactlyAt(start, end firefly.Point, c firefly.Color) {
	if start.X < 0 && end.X < 0 {
		return
	}
	if start.Y < 0 && end.Y < 0 {
		return
	}
	firefly.DrawLine(
		start, end,
		firefly.L(c, snakeWidth),
	)
	firefly.DrawCircle(
		firefly.Point{
			X: end.X - snakeWidth/2,
			Y: end.Y - snakeWidth/2,
		},
		snakeWidth,
		firefly.Solid(c),
	)
}
