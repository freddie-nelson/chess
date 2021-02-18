package main

import (
	"github.com/containerd/console"
)

// Size is the width and height of the board
const Size int = 8

// TermSize is the number of columns and rows the board occupies in the terminal
const TermSize int = 32

func main() {
	// TermSize = getTermSize()

	board := SetupInitialBoard()
	Draw(board)
}

func getTermSize() int {
	current := console.Current()

	ws, err := current.Size()
	if err != nil {
		return 10
	}

	return int(ws.Height)
}
