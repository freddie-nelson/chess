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
	output := GameState.board.ToString()
	fmt.Printf("\033[0;0H%s%v", output, time.Now())
}

// Tick reacts to changes in the game's state every tick
func (b *BoardEntity) Tick(e tl.Event) {
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
}

// SetupBoardLevel sets up the board level and returns it
func SetupBoardLevel() *tl.BaseLevel {
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
		Fg: tl.ColorWhite,
	})

	level.AddEntity(&BoardEntity{tl.NewEntity(0, 0, 0, 0)})

	return level
}