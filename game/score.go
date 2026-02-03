package game

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

// For how long (in frames) the snake is invincible after a collision.
const iFrames = 60

// Badges.
const (
	badgeBiteSelf     firefly.Badge = 1
	badgeBiteOther    firefly.Badge = 2
	badgeEat100Apples firefly.Badge = 3
)

// How long (in frames) the snake can go without food.
var hungerPeriod uint16

type Score struct {
	peer firefly.Peer

	// The current score. Cannot go below zero.
	val int16

	// Invisibility frames.
	// For how many frames from now the snake is invinsible.
	iframes uint8

	// How many more frames the snake can last without food.
	// If reaches zero, the scroe decrements by one step.
	hunger uint16

	ttl uint8

	color firefly.Color
}

func newScore(peer firefly.Peer) *Score {
	return &Score{
		peer:    peer,
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
		//
		// Hunger has no effect when the score is zero,
		// which prevents it from being activated at the beginning,
		// before the snake eats the first apple.
		if s.val != 0 {
			s.dec()
			s.hunger = hungerPeriod
			if s.val == 0 {
				snakes.deletePeer(s.peer)
				if s.peer == me {
					setTitle("ur snek ded cuz its hungie :(")
				} else {
					setTitle("aze snek got hungie, u win")
				}
			}
		}
	} else {
		s.hunger--
	}
	if s.ttl != 0 {
		s.ttl--
	}
}

// Increase the score.
//
// Triggered by [Snake] when eating an apple.
func (s *Score) inc() {
	if hungerPeriod > 10 {
		hungerPeriod -= 1
	}
	s.hunger = hungerPeriod
	s.val += 1
	firefly.AddProgress(s.peer, badgeEat100Apples, 1)
	s.color = firefly.ColorDarkGreen
	s.ttl = 60
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
	s.color = firefly.ColorRed
	s.ttl = 60
}
