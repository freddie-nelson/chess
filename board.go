package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Board stores details about the boards state
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

	b.selectedSpot = &board[0][7]
	b.selectedSpot.selected = true
	b.grid = &board

	// generate starting position
	startingFEN := "R6k/8/8/8/8/8/8/1Q6 w KQkq - 0 1"
	b.GenerateFromFENString(startingFEN)
}

// GenerateFromFENString creates a particular board position from a provided valid FEN string
func (b *Board) GenerateFromFENString(fen string) {
	piecePlacements := strings.Split(fen, "/")

	last := strings.Split(piecePlacements[7], " ")
	piecePlacements[7] = last[0]
	fields := last[1:]

	// current turn
	if fields[0] == "b" {
		GameState.turn = Black
	} else {
		GameState.turn = White
	}

	// castling rights
	GameState.blackCastling = &CastlingRights{false, false}
	GameState.whiteCastling = &CastlingRights{false, false}

	for _, rights := range fields[1] {
		if unicode.IsLower(rights) {
			if rights == 'k' {
				GameState.blackCastling.kingside = true
			} else {
				GameState.blackCastling.queenside = true
			}
		} else {
			if rights == 'K' {
				GameState.whiteCastling.kingside = true
			} else {
				GameState.whiteCastling.queenside = true
			}
		}
	}

	// en passant targets
	if fields[2] != "-" {
		file, rank := b.locationToFileAndRank(fields[2])
		b.grid[file][rank].passantTarget = 2
	}

	// fullmoves and halfmoves
	GameState.halfmoves = int(fields[3][0] - '0')
	GameState.fullmoves = int(fields[4][0] - '0')

	// place pieces
	for rank, fenRank := range piecePlacements {
		file := 0

		for _, char := range fenRank {
			var color int
			var class int

			if unicode.IsNumber(char) {
				file += int(char - '0')
				continue
			} else if unicode.IsLower(char) {
				color = Black
			} else {
				color = White
			}

			switch unicode.ToUpper(char) {
			case 'Q':
				class = Queen
			case 'K':
				class = King
			case 'R':
				class = Rook
			case 'B':
				class = Bishop
			case 'N':
				class = Knight
			case 'P':
				class = Pawn
			}

			b.grid[file][rank].containsPiece = true
			b.grid[file][rank].piece = &Piece{color: color, class: class}
			file++
		}
	}
}

func (b *Board) locationToFileAndRank(loc string) (int, int) {
	file := int(loc[0] - 'a')
	rank := 8 - int(loc[1]-'0')
	fmt.Printf(" file: %v rank: %v", file, rank)
	return file, rank
}

// ClearHighlighted sets all highlighted squares back to not highlighted
func (b *Board) ClearHighlighted() {
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			b.grid[file][rank].highlighted = false
			b.grid[file][rank].passantTarget--
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
	if (!b.selectedSpot.containsPiece || b.selectedSpot.piece.color == Black) && !b.selectedSpot.highlighted {
		return
	}

	if b.pickedSpot != nil {
		b.pickedSpot.picked = false
	}

	if b.selectedSpot.highlighted {
		b.MovePiece(b.pickedSpot, b.selectedSpot)
		b.pickedSpot = nil
		return
	}

	b.pickedSpot = b.selectedSpot
	b.pickedSpot.picked = true

	b.pickedSpot.piece.FindValidMoves(b.grid, b.pickedSpot.file, b.pickedSpot.rank, Black)
}

// IsSpotOffBoard returns true if the spot is not on the board
func (b *Board) IsSpotOffBoard(file int, rank int) bool {
	return file < 0 || file > Size-1 || rank < 0 || rank > Size-1
}

// MovePiece moves a piece from a start position on the board to the destination
func (b *Board) MovePiece(start *Spot, destination *Spot) {
	turnSuccessful := true

	piece := start.piece
	destinationPiece := destination.piece
	piece.moves++

	start.piece = nil
	start.containsPiece = false

	destination.piece = piece
	destination.containsPiece = true

	// if it is pawns first move and moved 2 places make spot behind pawn en passant target for next turn
	if destination.piece.class == Pawn && destination.piece.moves == 1 && destination.rank == 4 {
		b.grid[destination.file][destination.rank+1].passantTarget = 2
	}

	// if pawn move results in en passant, take piece behind destination
	var passantPiece *Piece
	if destination.passantTarget > 0 {
		passantSpot := &b.grid[destination.file][start.rank]
		passantPiece = passantSpot.piece
		passantSpot.piece = nil
		passantSpot.containsPiece = false
	}

	// if move puts player's king in check then revert the move
	if b.IsKingInCheck() {
		piece.moves--

		start.piece = piece
		start.containsPiece = true

		destination.piece = destinationPiece
		destination.containsPiece = destinationPiece != nil

		if destination.passantTarget > 0 {
			passantSpot := &b.grid[destination.file][start.rank]
			passantSpot.piece = passantPiece
			passantSpot.containsPiece = true
		}

		turnSuccessful = false

		if b.IsCheckmate() {
			GameState.ended = true
		}
	}

	// if turn was successfully played then pass turn to opponent
	if turnSuccessful {
		GameState.halfmoves++

		if GameState.turn == Black {
			GameState.fullmoves++
			GameState.turn = White
		} else {
			GameState.turn = Black
		}
	}

	if GameState.fullmoves == 50 {
		GameState.ended = true
	}

	// clear highlighted possible moves once piece has moved
	b.ClearHighlighted()
}

// IsKingInCheck goes through each opponent piece on the board and checks if they are attacking
// the player's king
// returns either true (the king is in check) or false (the king is not in check)
func (b *Board) IsKingInCheck() bool {

	// check if any opponent's piece puts the king in check
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			if b.grid[file][rank].containsPiece && b.grid[file][rank].piece.color == Black {
				_, inCheck := b.grid[file][rank].piece.FindValidMoves(b.grid, file, rank, White)
				if inCheck {
					// fmt.Printf(" file: %v rank: %v", file, rank)
					return true
				}
			}
		}
	}

	return false
}

// IsCheckmate checks if player is in checkmate
func (b *Board) IsCheckmate() bool {
	// find king on board
	var king *Spot
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			if b.grid[file][rank].containsPiece && b.grid[file][rank].piece.class == King && b.grid[file][rank].piece.color == White {
				king = &b.grid[file][rank]
			}
		}
	}

	// get king moves
	kingMoves, _ := king.piece.FindValidMoves(b.grid, king.file, king.rank, Black)

	// check if any opponent piece is attacking a spot the king can move to
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			if b.grid[file][rank].containsPiece && b.grid[file][rank].piece.color == Black {
				moves, _ := b.grid[file][rank].piece.FindValidMoves(b.grid, file, rank, Black)

				// check if the piece is attacking any of the kings moves
				for _, move := range moves {
					for i, kingMove := range kingMoves {
						if move.file == kingMove.file && move.rank == kingMove.rank {
							kingMoves = append(kingMoves[:], kingMoves[i+1:]...)
						}
					}
				}
			}
		}

		if len(kingMoves) == 0 {
			return true
		}

		return false
	}

	return false
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
