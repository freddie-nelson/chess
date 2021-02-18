package main

import (
	"fmt"
	"strings"
)

// SetupInitialBoard creates the initial chess board
func SetupInitialBoard() *[Size][Size]string {
	board := [Size][Size]string{}

	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			board[file][rank] = "P"
		}
	}

	return &board
}

// Draw draws the current board to the terminal
func Draw(board *[Size][Size]string) {
	output := ""

	border := strings.Repeat("||-------", Size) + "||\n"
	gap := strings.Repeat("||       ", Size) + "||\n"

	for rank := 0; rank < Size; rank++ {
		lines := [4]string{}

		lines[0] = border
		lines[1] = gap
		lines[3] = gap

		for file := 0; file < Size; file++ {
			piece := fmt.Sprintf("||   %s   ", board[file][rank])
			lines[2] += piece
		}

		lines[2] += "||\n"

		output += strings.Join(lines[:], "")
	}

	output += border

	fmt.Print(output)
}
