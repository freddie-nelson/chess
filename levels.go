package main

import (
	"fmt"
	"time"

	tl "github.com/JoelOtter/termloop"
)

// BoardEntity represents the board in the game space
type BoardEntity struct {
	*tl.Entity
}

// Draw draws the boards current state to the console
func (b *BoardEntity) Draw(s *tl.Screen) {
	// print board
	output := GameState.board.ToString()
	fmt.Printf("\033[4;0H%s", output)
}

// Tick reacts to changes in the game's state every tick
func (b *BoardEntity) Tick(e tl.Event) {
	// calculate deltaTime
	now := int(time.Now().UnixNano() / 1000000)
	if GameState.timeOfLastTick == 0 {
		GameState.timeOfLastTick = now
	}

	GameState.deltaTime = now - GameState.timeOfLastTick
	GameState.timeOfLastTick = now

	board := GameState.board

	if e.Type == tl.EventKey {
		switch e.Key {
		case tl.KeyArrowRight:
			board.ChangeSelectedSpot(1, 0)
		case tl.KeyArrowLeft:
			board.ChangeSelectedSpot(-1, 0)
		case tl.KeyArrowUp:
			board.ChangeSelectedSpot(0, -1)
		case tl.KeyArrowDown:
			board.ChangeSelectedSpot(0, 1)
		case tl.KeyEnter:
			board.PickSpot()
		}
	}

	// if GameState.turn == GameState.color {
	// 	GameState.time -= deltaTime
	// }
}

// SetupBoardLevel sets up the board level and returns it
func SetupBoardLevel() *tl.BaseLevel {
	level := tl.NewBaseLevel(tl.Cell{})

	level.AddEntity(&BoardEntity{tl.NewEntity(0, 0, 0, 0)})

	return level
}
