package main

import (
	tl "github.com/JoelOtter/termloop"
)

// Size is the width and height of the board
const Size int = 8

// GameState global game state
var GameState Game

func main() {
	game := tl.NewGame()
	screen := game.Screen()

	board := SetupInitialBoard()
	GameState.board = board

	boardLevel := SetupBoardLevel()
	screen.SetLevel(boardLevel)

	screen.SetFps(24)

	game.Start()
}
