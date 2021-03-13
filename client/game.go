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

func (g *Game) NextTurn(color int, opponentColor int) {
	// check for winning conditions
	if g.board.IsStalemate(opponentColor, color) {
		GameState.ended = true

		if g.board.IsKingInCheck(opponentColor, color, nil) {
			GameState.endState = "checkmate"
		} else {
			GameState.endState = "stalemate"
		}
	}

	GameState.halfmoves++
	if GameState.turn == Black {
		GameState.fullmoves++
		GameState.turn = White
	} else {
		GameState.turn = Black
	}

	if GameState.fullmoves == 50 {
		GameState.ended = true
	}
}
