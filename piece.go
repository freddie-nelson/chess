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
var PieceStrings []string = []string{"Q", "K", "R", "B", "Kn", "P"}

// Piece : generic class for a chess piece
type Piece struct {
	color int
	class int
	moves int
}

// FindValidMoves finds and returns all the legal moves a piece can make from it's current position
func (p *Piece) FindValidMoves(b *[Size][Size]Spot, file int, rank int) [Size][Size]bool {
	validMoves := [Size][Size]bool{}

	// bishop offsets
	bishopXOffs := []int{1, -1}
	bishopYOffs := []int{-1, 1}

	// rook offsets
	rookXOffs := []int{0, 0}
	rookYOffs := []int{-1, 1}

	// queen offsets
	queenOffs := []int{0, 0, -1, 1}

	GameState.board.ClearHighlighted()

	switch p.class {
	case Queen:
		calculateMovesFromOffsets(b, &validMoves, file, rank, queenOffs, queenOffs, Size, true)
	case King:
		calculateMovesFromOffsets(b, &validMoves, file, rank, queenOffs, queenOffs, 1, true)
	case Rook:
		calculateMovesFromOffsets(b, &validMoves, file, rank, rookXOffs, rookYOffs, Size, true)
		calculateMovesFromOffsets(b, &validMoves, file, rank, rookYOffs, rookXOffs, Size, true)
	case Bishop:
		calculateMovesFromOffsets(b, &validMoves, file, rank, bishopXOffs, bishopYOffs, Size, true)
	case Knight:
	case Pawn:
		if p.moves == 0 {
			calculateMovesFromOffsets(b, &validMoves, file, rank, []int{0}, []int{-1}, 2, false)
		} else {
			calculateMovesFromOffsets(b, &validMoves, file, rank, []int{0}, []int{-1}, 1, false)
		}

		checkIfPawnCanTake(b, &validMoves, file, rank)
	}

	return validMoves
}

func calculateMovesFromOffsets(b *[Size][Size]Spot, validMoves *[Size][Size]bool, file int, rank int, xOffs []int, yOffs []int, stopAfter int, canTake bool) {
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
					if spot.piece.color == Black && canTake {
						if spot.piece.class == King {
							break
						}

						spot.highlighted = true
						validMoves[currentFile][currentRank] = true
					}

					break
				} else {
					spot.highlighted = true
					validMoves[currentFile][currentRank] = true
				}
			}
		}
	}
}

func checkIfPawnCanTake(b *[Size][Size]Spot, validMoves *[Size][Size]bool, file int, rank int) {
	// calculate positions on board
	lFile := file - 1
	rFile := file + 1
	nextRank := rank - 1

	// left file
	if !GameState.board.IsSpotOffBoard(lFile, nextRank) {

		// can pawn take diagonally
		if b[lFile][nextRank].containsPiece && b[lFile][nextRank].piece.color == Black {
			b[lFile][nextRank].highlighted = true
			validMoves[lFile][nextRank] = true
		} else if b[lFile][rank].containsPiece && b[lFile][rank].piece.color == Black { // can pawn take en passant
			b[lFile][nextRank].highlighted = true
			b[lFile][nextRank].passantMove = true
			validMoves[lFile][nextRank] = true
		}
	}

	// right file
	if !GameState.board.IsSpotOffBoard(rFile, nextRank) {

		// can pawn take diagonally
		if b[rFile][nextRank].containsPiece && b[rFile][nextRank].piece.color == Black {
			b[rFile][nextRank].highlighted = true
			validMoves[rFile][nextRank] = true
		} else if b[rFile][rank].containsPiece && b[rFile][rank].piece.color == Black { // can pawn take en passant
			b[rFile][nextRank].highlighted = true
			b[rFile][nextRank].passantMove = true
			validMoves[rFile][nextRank] = true
		}
	}
}
