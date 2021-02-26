package main

// Enum type of piece
const (
	Queen int = iota
	King
	Rook
	Bishop
	Knight
	Pawn
)

// Enum color of piece
const (
	Black int = iota
	White
)

// PieceStrings an array containing the string values for outputting pieces to the screen
var PieceStrings []string = []string{"Q", "K", "R", "B", "N", "P"}

// Piece : generic class for a chess piece
type Piece struct {
	color int
	class int
	moves int
}

// FindValidMoves finds and returns all the legal moves a piece can make from it's current position
// returns array of all valid moves parallel to the board array
// and a boolean which represents if the piece is currently attacking the king
func (p *Piece) FindValidMoves(b *[Size][Size]Spot, file int, rank int, opponentColor int) ([]Spot, bool) {
	validMoves := make([]Spot, 0)

	// bishop offsets
	bishopXOffs := []int{1, -1}
	bishopYOffs := []int{-1, 1}

	// rook offsets
	rookXOffs := []int{0, 0}
	rookYOffs := []int{-1, 1}

	// queen offsets
	queenOffs := []int{0, 0, -1, 1}

	// knight offsets
	knightXOffs := []int{2, -2}
	knightYOffs := []int{1, -1}

	GameState.board.ClearHighlighted()

	checksKing := false

	switch p.class {
	case Queen:
		checksKing = calculateMovesFromOffsets(b, validMoves, file, rank, queenOffs, queenOffs, Size, true, opponentColor)
	case King:
		checksKing = calculateMovesFromOffsets(b, validMoves, file, rank, queenOffs, queenOffs, 1, true, opponentColor)
	case Rook:
		checksKing = calculateMovesFromOffsets(b, validMoves, file, rank, rookXOffs, rookYOffs, Size, true, opponentColor)
		if checksKing {
			calculateMovesFromOffsets(b, validMoves, file, rank, rookYOffs, rookXOffs, Size, true, opponentColor)
		} else {
			checksKing = calculateMovesFromOffsets(b, validMoves, file, rank, rookYOffs, rookXOffs, Size, true, opponentColor)
		}
	case Bishop:
		calculateMovesFromOffsets(b, validMoves, file, rank, bishopXOffs, bishopYOffs, Size, true, opponentColor)
	case Knight:
		calculateMovesFromOffsets(b, validMoves, file, rank, knightXOffs, knightYOffs, 1, true, opponentColor)
		calculateMovesFromOffsets(b, validMoves, file, rank, knightYOffs, knightXOffs, 1, true, opponentColor)
	case Pawn:
		if p.moves == 0 {
			calculateMovesFromOffsets(b, validMoves, file, rank, []int{0}, []int{-1}, 2, false, opponentColor)
		} else {
			calculateMovesFromOffsets(b, validMoves, file, rank, []int{0}, []int{-1}, 1, false, opponentColor)
		}

		checksKing = checkIfPawnCanTake(b, validMoves, file, rank, opponentColor)
	}

	return validMoves, checksKing
}

func calculateMovesFromOffsets(b *[Size][Size]Spot, validMoves []Spot, file int, rank int, xOffs []int, yOffs []int, stopAfter int, canTake bool, opponentColor int) bool {
	checksKing := false

	// use offsets to jump across board in the way the piece would
	for _, xOff := range xOffs {
		for _, yOff := range yOffs {

			// loop through each spot until end of board or stopAfter is reached
			// if we run into a spot that is occupied then break
			// highlight occupied spot before break if it is an opponent's piece
			for i := 1; i < Size && i <= stopAfter; i++ {
				currentFile := file + xOff*i
				currentRank := rank + yOff*i
				if GameState.board.IsSpotOffBoard(currentFile, currentRank) {
					break
				}

				spot := &b[currentFile][currentRank]
				if spot.containsPiece {
					if spot.piece.color == opponentColor && canTake {
						if spot.piece.class == King {
							checksKing = true
							break
						}

						spot.highlighted = true
						validMoves = append(validMoves, Spot{file: currentFile, rank: currentRank})
					}

					break
				} else {
					spot.highlighted = true
					validMoves = append(validMoves, Spot{file: currentFile, rank: currentRank})
				}
			}
		}
	}

	return checksKing
}

func checkIfPawnCanTake(b *[Size][Size]Spot, validMoves []Spot, file int, rank int, opponentColor int) bool {
	// calculate positions on board
	lFile := file - 1
	rFile := file + 1
	nextRank := rank - 1

	checksKing := false

	// left file
	if !GameState.board.IsSpotOffBoard(lFile, nextRank) {

		// can pawn take diagonally
		if b[lFile][nextRank].containsPiece && b[lFile][nextRank].piece.color == opponentColor {
			if b[lFile][nextRank].piece.class == King {
				checksKing = true
			} else {
				b[lFile][nextRank].highlighted = true
				validMoves = append(validMoves, Spot{file: lFile, rank: nextRank})
			}
		} else if !b[lFile][nextRank].containsPiece && b[lFile][nextRank].passantTarget > 0 { // can pawn take en passant
			b[lFile][nextRank].highlighted = true
			validMoves = append(validMoves, Spot{file: lFile, rank: nextRank})
		}
	}

	// right file
	if !GameState.board.IsSpotOffBoard(rFile, nextRank) {

		// can pawn take diagonally
		if b[rFile][nextRank].containsPiece && b[rFile][nextRank].piece.color == opponentColor {
			if b[rFile][nextRank].piece.class == King {
				checksKing = true
			} else {
				b[rFile][nextRank].highlighted = true
				validMoves = append(validMoves, Spot{file: rFile, rank: nextRank})
			}
		} else if !b[rFile][nextRank].containsPiece && b[rFile][nextRank].passantTarget > 0 { // can pawn take en passant
			b[rFile][nextRank].highlighted = true
			validMoves = append(validMoves, Spot{file: rFile, rank: nextRank})
		}
	}

	return checksKing
}
