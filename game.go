package main

// CastlingRights stores what side a player can castle
type CastlingRights struct {
	queenside bool
	kingside  bool
}

// Game stores values about the games current state
type Game struct {
	board         *Board
	ended         bool
	turn          int
	halfmoves     int
	fullmoves     int
	whiteCastling *CastlingRights
	blackCastling *CastlingRights
}
