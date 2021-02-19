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
		calculateMovesFromOffsets(b, &validMoves, file, rank, queenOffs, queenOffs, Size)
	case King:
		calculateMovesFromOffsets(b, &validMoves, file, rank, queenOffs, queenOffs, 1)
	case Rook:
		calculateMovesFromOffsets(b, &validMoves, file, rank, rookXOffs, rookYOffs, Size)
		calculateMovesFromOffsets(b, &validMoves, file, rank, rookYOffs, rookXOffs, Size)
	case Bishop:
		calculateMovesFromOffsets(b, &validMoves, file, rank, bishopXOffs, bishopYOffs, Size)
	case Knight:
	case Pawn:
		if p.moves == 0 {
			calculateMovesFromOffsets(b, &validMoves, file, rank, []int{0}, []int{-1}, 2)
		} else {
			calculateMovesFromOffsets(b, &validMoves, file, rank, []int{0}, []int{-1}, 1)

		}
	}

	return validMoves
}

func calculateMovesFromOffsets(board *[Size][Size]Spot, validMoves *[Size][Size]bool, file int, rank int, xOffs []int, yOffs []int, stopAfter int) {
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

				spot := &board[currentFile][currentRank]
				if spot.containsPiece {
					if spot.piece.color == Black {
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
