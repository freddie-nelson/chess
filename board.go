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
			board[file][rank] = " "
		}
	}

	board[0][0] = "A"
	board[1][1] = "B"
	board[2][2] = "C"
	board[3][3] = "D"
	board[4][4] = "E"

	return &board
}

// Draw draws the current board to the terminal
func Draw(board *[Size][Size]string) {
	output := ""

	resetColor := "\033[0m"
	darkSquareColor := "\033[43m"
	lightSquareColor := "\033[100m"
	gapChar := " "
	pieceLine := 2
	spotSize := 11

	for rank := 0; rank < Size; rank++ {
		lines := [5]string{}

		for file := 0; file < Size; file++ {
			bgColor := darkSquareColor
			if (file+rank)%2 == 0 {
				bgColor = lightSquareColor
			}

			gap := bgColor + strings.Repeat(gapChar, spotSize) + resetColor

			for i := 0; i < len(lines); i++ {
				if i == pieceLine {
					margin := bgColor + strings.Repeat(gapChar, spotSize/2)
					piece := fmt.Sprintf("%s%s%s%s", margin, board[file][rank], margin, resetColor)
					lines[i] += piece
					continue
				}

				lines[i] += gap
			}

		}

		for i := 0; i < len(lines); i++ {
			lines[i] += "\n"
		}

		output += strings.Join(lines[:], "")
	}

	fmt.Print(output)
}
