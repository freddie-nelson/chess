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
}
