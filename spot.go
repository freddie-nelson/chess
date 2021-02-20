package main

// Spot identifies a location on the board
type Spot struct {
	piece         *Piece
	containsPiece bool
	file          int
	rank          int
	selected      bool
	picked        bool
	highlighted   bool
	passantMove   bool
}
