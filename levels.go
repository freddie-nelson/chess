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
	output := BoardToString(GameState.board)
	fmt.Printf("\033[0;0H%s%v", output, time.Now())
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
