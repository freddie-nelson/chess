package main

// Spot identifies a location on the board
type Spot struct {
	piece         *Piece `json: "omitempty"`
	containsPiece bool
	file          int
	rank          int
}
