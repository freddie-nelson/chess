package main

import (
	"github.com/containerd/console"
)

// Size is the width and height of the board
const Size int = 8

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
