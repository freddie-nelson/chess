package main

// type of piece
const (
	Queen int = iota
	King
	Rook
	Bishop
	Knight
	Pawn
)

// color of piece
const (
	Black int = iota
	White
)

// Piece : generic class for a chess piece
type Piece struct {
	color int
	class int
	file  int
	rank  int
}
