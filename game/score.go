package game

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

const (
	// How long (in frames) the snake can go without food.
	hungerPeriod = 6 * 60

	// For how long (in frames) the snake is invulnerable after a collision.
	iFrames = 60

	badgeBiteSelf     firefly.Badge = 1
	badgeBiteOther    firefly.Badge = 2
	badgeEat100Apples firefly.Badge = 3
)

type Score struct {
	// The current score.
	// Cannot go below zero.
	val int

	// Invisibility frames.
	// For how many frames from now the snake is invinsible.
	iframes int

	// How many more frames the snake can last without food.
	// If reaches zero, the scroe decrements by one step.
	hunger int
}

func newScore() *Score {
	return &Score{
		hunger:  hungerPeriod,
		iframes: iFrames,
	}
}

// update the score.
//
// Checks for collisions and iframes and decrements the score if needed.
func (s *Score) update() {
	if s.iframes > 0 {
		s.iframes -= 1
	}
	if s.hunger == 0 {
		// Hungry. Decrese the score and start counting again.
		// s.dec()
		s.hunger = hungerPeriod
	} else {
		s.hunger -= 1
	}
}

// Increase the score.
//
// Triggered by [Snake] when eating an apple.
func (s *Score) inc() {
	s.hunger = hungerPeriod
	s.val += 1
	firefly.AddProgress(firefly.Combined, badgeEat100Apples, 1)
}

// Decrease the score.
//
// Triggered by the score itself when the snake collides with itself.
func (s *Score) dec() {
	if s.iframes > 0 {
		return
	}
	s.iframes = iFrames
	if s.val > 0 {
		s.val -= (s.val/5 + 1)
	}
}
