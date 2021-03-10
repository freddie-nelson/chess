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

	whiteCastling *CastlingRights
	blackCastling *CastlingRights

	you      *User
	opponent *User

	timeOfLastTick int
	deltaTime      int
}
