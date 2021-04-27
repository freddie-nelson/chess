package main

// CastlingRights stores what side a player can castle
type CastlingRights struct {
	queenside bool
	kingside  bool
}

// GameController controls top level game logic and handles server connections
type GameController struct {
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

// NextTurn performs end game state checks and if game does not end then proceeds to next turn
func (g *GameController) NextTurn(color int, opponentColor int) {
	// check for winning conditions
	if g.board.IsStalemate(opponentColor, color) {
		Game.ended = true

		if g.board.IsKingInCheck(opponentColor, color, nil) {
			Game.endState = "checkmate"
		} else {
			Game.endState = "stalemate"
		}
	}

	Game.halfmoves++
	if Game.turn == Black {
		Game.fullmoves++
		Game.turn = White
	} else {
		Game.turn = Black
	}

	if Game.fullmoves == 50 {
		Game.ended = true
	}
}
