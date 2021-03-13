package main

import (
	tl "github.com/JoelOtter/termloop"
)

// Size is the width and height of the board
const Size int = 8

// Game global game state
var Game GameController

func main() {
	game := tl.NewGame()
	screen := game.Screen()

	board := Board{}
	board.Setup()
	Game.board = &board
	Game.color = White
	Game.opponentColor = Black

	// setup users temp
	Game.you = &User{"Freddie", 600000, false}
	Game.opponent = &User{"GM Hikaru", 600000, true}

	mainMenuLevel := SetupMainMenuLevel()
	screen.SetLevel(mainMenuLevel)

	screen.SetFps(24)

	game.Start()
}
