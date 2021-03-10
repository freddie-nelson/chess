package main

// CastlingRights stores what side a player can castle
type CastlingRights struct {
	queenside bool
	kingside  bool
}

// Game stores values about the games current state
type Game struct {
	color         int
	opponentColor int
	board         *Board
	ended         bool
	endState      string
	turn          int
	halfmoves     int
	fullmoves     int
	time          int
	opponentTime  int
	whiteCastling *CastlingRights
	blackCastling *CastlingRights

	timeOfLastTick int
	deltaTime      int
}
