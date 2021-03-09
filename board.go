package main

import (
	"fmt"
	"strings"
	"time"
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
	startingFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
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
	if ((!b.selectedSpot.containsPiece || b.selectedSpot.piece.color != GameState.color) && !b.selectedSpot.highlighted) || GameState.turn != GameState.color {
		if !b.selectedSpot.containsPiece {
			if b.pickedSpot != nil {
				b.pickedSpot.picked = false
			}

			b.pickedSpot = nil
			b.ClearHighlighted()
		}

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

	validMoves, _ := b.pickedSpot.piece.FindValidMoves(b.grid, b.pickedSpot.file, b.pickedSpot.rank, Black, true)
	b.highlightMoves(validMoves)
}

func (b *Board) highlightMoves(moves []Spot) {
	for _, m := range moves {
		b.grid[m.file][m.rank].highlighted = true
	}
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
	opponentColor := Black
	if GameState.turn == Black {
		opponentColor = White
	}
	if b.IsKingInCheck(GameState.turn, opponentColor, nil) {
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
	}

	// if turn was successfully played then pass turn to opponent
	if turnSuccessful {
		b.nextTurn(GameState.turn, opponentColor)
	}

	// clear highlighted possible moves once piece has moved
	b.ClearHighlighted()
}

func (b *Board) nextTurn(color int, opponentColor int) {
	// check for winning conditions
	if b.IsStalemate(opponentColor, color) {
		GameState.ended = true

		if b.IsKingInCheck(opponentColor, color, nil) {
			GameState.endState = "checkmate"
		} else {
			GameState.endState = "stalemate"
		}
	}

	GameState.halfmoves++
	if GameState.turn == Black {
		GameState.fullmoves++
		GameState.turn = White
	} else {
		GameState.turn = Black
	}

	if GameState.fullmoves == 50 {
		GameState.ended = true
	}
}

// IsKingInCheck goes through each opponent piece on the board and checks if they are attacking
// color's king
// returns either true (the king is in check) or false (the king is not in check)
func (b *Board) IsKingInCheck(color int, opponentColor int, simulatedBoard *[Size][Size]Spot) bool {
	board := b.grid
	if simulatedBoard != nil {
		board = simulatedBoard
	}

	// check if any opponent's piece puts the king in check
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			if board[file][rank].containsPiece && board[file][rank].piece.color == opponentColor {
				_, inCheck := board[file][rank].piece.FindValidMoves(board, file, rank, color, false)
				if inCheck {
					return true
				}
			}
		}
	}

	return false
}

// IsStalemate returns true if color cannot play any moves but is not in check
func (b *Board) IsStalemate(color int, opponentColor int) bool {
	king := b.GetKingSpot(color)
	kingMoves, _ := king.piece.FindValidMoves(b.grid, king.file, king.rank, opponentColor, true)

	// when king cannot move out of check
	// check if any move by color can get king out of check
	if len(kingMoves) == 0 {
		for rank := 0; rank < Size; rank++ {
			for file := 0; file < Size; file++ {
				if b.grid[file][rank].containsPiece && b.grid[file][rank].piece.color == color && b.grid[file][rank].piece.class != King {
					piece := b.grid[file][rank].piece
					moves, _ := piece.FindValidMoves(b.grid, file, rank, opponentColor, true)

					// since moves are pruned for illegal moves
					// if any move is available then it will put king out of check
					if len(moves) > 0 {
						return false
					}
				}
			}
		}
	} else {
		return false
	}

	return true
}

// GetKingSpot returns the spot that contains the king of color color
func (b *Board) GetKingSpot(color int) *Spot {
	var king *Spot
	for rank := 0; rank < Size; rank++ {
		for file := 0; file < Size; file++ {
			if b.grid[file][rank].containsPiece && b.grid[file][rank].piece.class == King && b.grid[file][rank].piece.color == color {
				king = &b.grid[file][rank]
			}
		}
	}

	return king
}

// ToString returns the board's current state as a single string
func (b *Board) ToString() string {
	output := ""

	// set bg back to default
	resetColor := "\033[0m"

	// square bg colors
	lightSquareColor := "\033[48;2;240;217;181m"
	darkSquareColor := "\033[48;2;181;136;99m"
	selectedLightSquareColor := "\033[48;2;205;210;106m"
	selectedDarkSquareColor := "\033[48;2;170;162;58m"
	highlightedLightSquareColor := "\033[48;2;183;176;170m"
	highlightedDarkSquareColor := "\033[48;2;138;114;107m"
	pickedLightSquareColor := "\033[48;2;109;159;88m"
	pickedDarkSquareColor := "\033[48;2;85;126;56m"

	// piece colors
	blackPieceColor := "\033[38;2;0;0;0m"
	whitePieceColor := "\033[38;2;255;255;255m"

	// grid coords colors
	darkCoordColor := "\033[38;2;240;217;181m"
	lightCoordColor := "\033[38;2;181;136;99m"

	// characters
	gapChar := " "

	pieceLine := 1
	const spotCols int = 7
	const spotRows int = 3

	for rank := 0; rank < Size; rank++ {
		lines := [spotRows]string{}

		for file := 0; file < Size; file++ {
			// select colour for checkered pattern
			bgColor := darkSquareColor
			if (file+rank)%2 == 0 {
				bgColor = lightSquareColor
			}

			spot := b.grid[file][rank]
			if spot.picked {
				if bgColor == darkSquareColor {
					bgColor = pickedDarkSquareColor
				} else {
					bgColor = pickedLightSquareColor
				}
			} else if spot.selected {
				if bgColor == darkSquareColor {
					bgColor = selectedDarkSquareColor
				} else {
					bgColor = selectedLightSquareColor
				}
			} else if spot.highlighted {
				if bgColor == darkSquareColor {
					bgColor = highlightedDarkSquareColor
				} else {
					bgColor = highlightedLightSquareColor
				}
			}

			gap := bgColor + strings.Repeat(gapChar, spotCols) + resetColor

			for i := 0; i < len(lines); i++ {
				if i == pieceLine {
					margin := bgColor + strings.Repeat(gapChar, spotCols/2)
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
				} else if i == 0 && file == 0 {
					textColor := ""
					if bgColor == darkSquareColor {
						textColor = darkCoordColor
					} else {
						textColor = lightCoordColor
					}

					lines[i] += bgColor + textColor + fmt.Sprint(Size-rank) + strings.Repeat(gapChar, spotCols-1) + resetColor
				} else if i == len(lines)-1 && rank == Size-1 {
					textColor := ""
					if bgColor == darkSquareColor {
						textColor = darkCoordColor
					} else {
						textColor = lightCoordColor
					}

					lines[i] += bgColor + strings.Repeat(gapChar, spotCols-1) + textColor + string(rune(file+97)) + resetColor
				} else {
					lines[i] += gap
				}
			}

		}

		for i := 0; i < len(lines); i++ {
			lines[i] = lines[i] + "\n"
			output += lines[i]
		}
	}

	// add timers
	top := b.createTimerString(GameState.opponentTime, false, spotCols, spotRows, resetColor)
	bottom := b.createTimerString(GameState.time, true, spotCols, spotRows, resetColor)

	output += top + bottom

	return output
}

func (b *Board) createTimerString(timerMs int, bottom bool, spotCols int, spotRows int, resetColor string) string {
	timer := ""
	timerCols := 10
	timerLines := 3
	timeLine := 1
	timerBgColor := "\033[48;2;0;0;0m"
	timerTextColor := "\033[38;2;255;255;255m"

	line := 0

	for i := 0; i < timerLines; i++ {
		line++
		if bottom {
			line = (spotRows * Size) - i
		}

		if i == timeLine {
			t, _ := time.ParseDuration(fmt.Sprintf("%vms", timerMs))
			mins := int(t.Minutes())
			secs := int(t.Seconds()) - mins*60
			ms := int(t.Milliseconds()) - (mins * 60 * 1000) - (secs * 1000)

			timeString := fmt.Sprintf("%v:%02v.%v", mins, secs, string(fmt.Sprint(ms)[0]))
			gapCount := timerCols - len(timeString) - 1

			timer += fmt.Sprintf("\033[%v;%vH%s%s%s%s %s", line, (spotCols*Size)+1, timerBgColor, timerTextColor, strings.Repeat(" ", gapCount), timeString, resetColor)
		} else {
			timer += fmt.Sprintf("\033[%v;%vH%s%s%s%s", line, (spotCols*Size)+1, timerBgColor, timerTextColor, strings.Repeat(" ", timerCols), resetColor)
		}
	}

	return timer
}
