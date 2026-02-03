package game

import "github.com/firefly-zero/firefly-go/firefly"

const (
	singleplayer firefly.Board = 1
	multiplayer  firefly.Board = 2
)

func updateLeaderBoard() {
	board := singleplayer
	if isMultiplayer {
		board = multiplayer
	}
	for _, snake := range snakes.items {
		if snake.score.val != 0 {
			firefly.AddScore(snake.peer, board, snake.score.val)
		}
	}
}
