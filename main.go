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

	board := Board{}
	board.Setup()
	GameState.board = &board
	GameState.color = White
	GameState.opponentColor = Black

	// setup users temp
	GameState.you = &User{"Freddie", 600000, false}
	GameState.opponent = &User{"GM Hikaru", 600000, true}

	boardLevel := SetupBoardLevel()
	screen.SetLevel(boardLevel)

	screen.SetFps(24)

	game.Start()
}
