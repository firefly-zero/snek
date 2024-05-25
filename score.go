package main

import (
	"strconv"

	"github.com/firefly-zero/firefly-go/firefly"
)

// How long the snake can go without food.
const HungerPeriod = 5 * 60

// For how long the snake is invulnerable after a collision.
const IFrames = 60

var score Score

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

func NewScore() Score {
	return Score{
		hunger:  HungerPeriod,
		iframes: IFrames,
	}
}

// Update the score.
//
// Checks for collisions and iframes and decrements the score if needed.
func (s *Score) Update(snake *Snake) {
	if s.iframes > 0 {
		s.iframes -= 1
	}
	if s.hunger == 0 {
		// Hungry. Decrese the score and start counting again.
		s.dec()
		s.hunger = HungerPeriod
	} else {
		s.hunger -= 1
	}
	if snake.Collides(snake.Mouth) {
		score.dec()
	}
}

// Increase the score.
//
// Triggered by [Snake] when eating an apple.
func (s *Score) Inc() {
	s.hunger = HungerPeriod
	s.val += 1
}

// Decreas the score.
//
// Triggered by the score itself when the snake collides with itself.
func (s *Score) dec() {
	if s.iframes > 0 {
		return
	}
	s.iframes = IFrames
	if s.val > 0 {
		s.val -= (s.val/5 + 1)
	}
}

// Show the score in the corner of the screen.
func (s Score) Render() {
	firefly.DrawText(
		strconv.Itoa(s.val), font,
		firefly.Point{X: 10, Y: 10},
		firefly.ColorDarkBlue,
	)
}
