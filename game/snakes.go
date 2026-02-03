package game

import (
	"fmt"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Snakes struct {
	items []*Snake
}

func newSnakes() *Snakes {
	peers := firefly.GetPeers().Slice()

	// Set the global hunger period based on the number of players.
	// The more people play, the longer it takes for one snake
	// to get an apple (because of competition).
	hungerPeriodSeconds := 4 + len(peers)
	hungerPeriod = uint16(hungerPeriodSeconds) * 60

	isMultiplayer := len(peers) != 1
	snakes := make([]*Snake, len(peers))
	for i, peer := range peers {
		snakes[i] = newSnake(peer, isMultiplayer)
	}
	return &Snakes{snakes}
}

func (ss *Snakes) update() {
	if ss == nil {
		return
	}
	var best *Snake
	var bestScore int16 = 0
	for _, snake := range ss.items {
		snake.crown = false
		snake.update()
		snake.tryEat()
		score := snake.score.val
		if score == bestScore {
			// Nobody's the best if there is a tie.
			best = nil
		} else if score > bestScore {
			best = snake
			bestScore = score
		}
	}
	// In multiplayer, render a crown on the best snake.
	if len(ss.items) > 1 && best != nil {
		best.crown = true
	}

	for i, s1 := range snakes.items {
		for j, s2 := range snakes.items {
			sameSnake := i == j
			if !s1.bites(sameSnake, s2) {
				continue
			}
			if sameSnake {
				firefly.AddProgress(s1.peer, badgeBiteSelf, 1)
			} else {
				firefly.AddProgress(s1.peer, badgeBiteOther, 1)
			}
			s1.eye.hurt = true
			s1.score.dec()

			// If the snake reached zero score, handle game over.
			if s1.score.val != 0 {
				continue
			}
			snakes.deleteSnake(s1)
			if sameSnake {
				if s1.peer == me {
					setTitle("u bit urself :(")
				} else {
					setTitle("other snek bit itself, u win")
				}
			} else {
				if s1.peer == me {
					setTitle("u lose :(")
				} else {
					setTitle("u win")
				}
			}
		}
	}
}

func (ss *Snakes) deleteSnake(tar *Snake) {
	newItems := make([]*Snake, 0, len(ss.items)-1)
	for _, s := range ss.items {
		if s != tar {
			newItems = append(newItems, s)
		}
	}
	ss.items = newItems
}

func (ss *Snakes) deletePeer(peer firefly.Peer) {
	newItems := make([]*Snake, 0, len(ss.items)-1)
	for _, s := range ss.items {
		if s.peer != peer {
			newItems = append(newItems, s)
		}
	}
	ss.items = newItems
}

// Check if the game over screen should be shown.
//
// In single-player, we end the game when the only snake dies.
// In multiplayer, we end the game when only one snake is left.
func (ss *Snakes) gameOver() bool {
	firefly.LogDebug(fmt.Sprintf("n snakes: %d", len(ss.items)))
	if len(ss.items) == 0 {
		return true
	}
	if len(ss.items) == 1 {
		return firefly.GetPeers().Len() > 1
	}
	return false
}

func (ss *Snakes) render() {
	if ss == nil {
		return
	}
	for _, snake := range ss.items {
		snake.render()
	}
}

// Check if an apple placed at the given point would collide with any snake.
//
// Used to pick a spot for a new apple position.
func (ss *Snakes) appleInside(pos firefly.Point) bool {
	if ss == nil {
		return false
	}
	for _, s := range ss.items {
		if s.appleCollides(pos) {
			return true
		}
	}
	return false
}
