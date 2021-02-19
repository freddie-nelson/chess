package main

import (
	"fmt"
	"strings"
)

// SetupInitialBoard creates the initial chess board
func SetupInitialBoard() *[Size][Size]Spot {
	board := [Size][Size]Spot{}

	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			board[file][rank] = Spot{file: file, rank: rank}
		}
	}

	// oppenent's pieces
	placeBackRank(&board, 0, Black)
	placePawnRank(&board, 1, Black)

	// player's pieces
	placeBackRank(&board, 7, White)
	placePawnRank(&board, 6, White)

	GameState.selectedSpot = &board[0][7]
	GameState.selectedSpot.selected = true

	return &board
}

func placeBackRank(board *[Size][Size]Spot, rank int, color int) {
	board[0][rank] = Spot{&Piece{color, Rook}, true, 0, rank, false}
	board[1][rank] = Spot{&Piece{color, Knight}, true, 1, rank, false}
	board[2][rank] = Spot{&Piece{color, Bishop}, true, 2, rank, false}
	board[3][rank] = Spot{&Piece{color, Queen}, true, 3, rank, false}
	board[4][rank] = Spot{&Piece{color, King}, true, 4, rank, false}
	board[5][rank] = Spot{&Piece{color, Bishop}, true, 5, rank, false}
	board[6][rank] = Spot{&Piece{color, Knight}, true, 6, rank, false}
	board[7][rank] = Spot{&Piece{color, Rook}, true, 7, rank, false}
}

func placePawnRank(board *[Size][Size]Spot, rank int, color int) {
	for file := 0; file < Size; file++ {
		board[file][rank] = Spot{&Piece{color, Pawn}, true, file, rank, false}
	}
}

// BoardToString returns the board's current state as a single string
func BoardToString(board *[Size][Size]Spot) string {
	output := ""

	resetColor := "\033[0m"

	// square bg colors
	darkSquareColor := "\033[100m"
	lightSquareColor := "\033[47m"
	selectedSquareColor := "\033[41m"

	// piece colors
	blackPieceColor := "\033[30m"
	whitePieceColor := "\033[37;1m"

	// characters
	gapChar := " "

	pieceLine := 2
	spotSize := 11

	for rank := 0; rank < Size; rank++ {
		lines := [5]string{}

		for file := 0; file < Size; file++ {
			// select colour for checkered pattern
			bgColor := darkSquareColor
			if (file+rank)%2 == 0 {
				bgColor = lightSquareColor
			}

			if board[file][rank].selected {
				bgColor = selectedSquareColor
			}

			gap := bgColor + strings.Repeat(gapChar, spotSize) + resetColor

			for i := 0; i < len(lines); i++ {
				if i == pieceLine {
					margin := bgColor + strings.Repeat(gapChar, spotSize/2)

					spot := board[file][rank]
					spotStr := " "

					if spot.containsPiece {
						piece := spot.piece
						pieceColor := whitePieceColor
						if piece.color == Black {
							pieceColor = blackPieceColor
						}

						pieceStr := PieceStrings[piece.class]

						if piece.class == Knight {
							pieceStr = "\b" + pieceStr
						}

						spotStr = pieceColor + pieceStr
					}

					line := fmt.Sprintf("%s%s%s%s", margin, spotStr, margin, resetColor)
					lines[i] += line
					continue
				}

				lines[i] += gap
			}

		}

		for i := 0; i < len(lines); i++ {
			lines[i] = "\r" + lines[i] + "\n"
			output += lines[i]
		}
	}

	return output
}
