package main

import (
	"snakeoop/internal/game"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			game.End("recover")
		}
	}()

	game.Start()
}
