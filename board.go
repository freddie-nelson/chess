package main

import (
	"fmt"
	"strings"
)

type Board struct {
	grid         *[Size][Size]Spot
	selectedSpot *Spot
	pickedSpot   *Spot
}

// Setup creates the initial chess board
func (b *Board) Setup() {
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

	b.selectedSpot = &board[0][7]
	b.selectedSpot.selected = true

	b.grid = &board
}

func placeBackRank(board *[Size][Size]Spot, rank int, color int) {
	board[0][rank] = Spot{&Piece{color: color, class: Rook}, true, 0, rank, false, false, false}
	board[1][rank] = Spot{&Piece{color: color, class: Knight}, true, 1, rank, false, false, false}
	board[2][rank] = Spot{&Piece{color: color, class: Bishop}, true, 2, rank, false, false, false}
	board[3][rank] = Spot{&Piece{color: color, class: Queen}, true, 3, rank, false, false, false}
	board[4][rank] = Spot{&Piece{color: color, class: King}, true, 4, rank, false, false, false}
	board[5][rank] = Spot{&Piece{color: color, class: Bishop}, true, 5, rank, false, false, false}
	board[6][rank] = Spot{&Piece{color: color, class: Knight}, true, 6, rank, false, false, false}
	board[7][rank] = Spot{&Piece{color: color, class: Rook}, true, 7, rank, false, false, false}
}

func placePawnRank(board *[Size][Size]Spot, rank int, color int) {
	for file := 0; file < Size; file++ {
		board[file][rank] = Spot{&Piece{color, Pawn, 0}, true, file, rank, false, false, false}
	}
}

// ClearHighlighted sets all highlighted squares back to not highlighted
func (b *Board) ClearHighlighted() {
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			b.grid[file][rank].highlighted = false
		}
	}
}

// SetSelectedSpot unselects the current spot and sets the new one
func (b *Board) SetSelectedSpot(file int, rank int) {
	// check if desired spot exists
	if b.IsSpotOffBoard(file, rank) {
		return
	}

	// if a spot is already selected, unselect it
	if b.selectedSpot != nil {
		b.selectedSpot.selected = false
	}

	// set new spot
	b.selectedSpot = &b.grid[file][rank]
	b.selectedSpot.selected = true
}

// ChangeSelectedSpot sets the new selected spot based on an offset from the current selected spot
func (b *Board) ChangeSelectedSpot(fileOff int, rankOff int) {
	current := b.selectedSpot
	if current == nil {
		return
	}

	file := current.file + fileOff
	rank := current.rank + rankOff

	b.SetSelectedSpot(file, rank)
}

// PickSpot picks the current selected spot
func (b *Board) PickSpot() {
	// prevent player from picking spots that don't contain a piece
	// or don't belong to them
	if !b.selectedSpot.containsPiece {
		return
	} else if b.selectedSpot.piece.color == Black {
		return
	}

	if b.pickedSpot != nil {
		b.pickedSpot.picked = false
	}

	b.pickedSpot = b.selectedSpot
	b.pickedSpot.picked = true

	b.pickedSpot.piece.FindValidMoves(b.grid, b.pickedSpot.file, b.pickedSpot.rank)
}

// IsSpotOffBoard returns true if the spot is not on the board
func (b *Board) IsSpotOffBoard(file int, rank int) bool {
	return file < 0 || file > Size-1 || rank < 0 || rank > Size-1
}

// ToString returns the board's current state as a single string
func (b *Board) ToString() string {
	output := ""

	resetColor := "\033[0m"

	// square bg colors
	darkSquareColor := "\033[100m"
	lightSquareColor := "\033[47m"
	selectedSquareColor := "\033[41m"
	pickedSquareColor := "\033[42m"
	highlightedSquareColor := "\033[43m"

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

			spot := b.grid[file][rank]
			if spot.picked {
				bgColor = pickedSquareColor
			} else if spot.selected {
				bgColor = selectedSquareColor
			} else if spot.highlighted {
				bgColor = highlightedSquareColor
			}

			gap := bgColor + strings.Repeat(gapChar, spotSize) + resetColor

			for i := 0; i < len(lines); i++ {
				if i == pieceLine {
					margin := bgColor + strings.Repeat(gapChar, spotSize/2)
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
			lines[i] = lines[i] + "\n"
			output += lines[i]
		}
	}

	return output
}
