package main

import (
	"fmt"
	"time"

	tl "github.com/JoelOtter/termloop"
)

// BoardEntity represents the board in the game space
type GameListener struct {
	*tl.Entity
}

// Draw draws the boards current state to the console
func (b *GameListener) Draw(s *tl.Screen) {
	// print board
	output := GameState.board.ToString()
	fmt.Printf("\033[4;0H%s", output)
}

// Tick reacts to changes in the game's state every tick
func (b *GameListener) Tick(e tl.Event) {
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

// SetupGameLevel sets up the game level and returns it
func SetupGameLevel() *tl.BaseLevel {
	level := tl.NewBaseLevel(tl.Cell{})

	level.AddEntity(&GameListener{tl.NewEntity(0, 0, 0, 0)})

	return level
}

/* MAIN MENU */

// MenuListener entity that listens for keyboard events on the menu
type MenuListener struct {
	*tl.Entity
	buttons     []*tl.Rectangle
	buttonsText []*tl.Text
	currentBtn  int
}

// Tick executes events every tick
func (ml *MenuListener) Tick(e tl.Event) {
	if e.Type == tl.EventKey {
		// remove highlight from current button
		if ml.currentBtn != 0 {
			ml.buttons[ml.currentBtn-1].SetColor(tl.ColorBlack)
			ml.buttonsText[ml.currentBtn-1].SetColor(tl.ColorWhite, tl.ColorBlack)
		}

		switch e.Key {
		case tl.KeyArrowDown:
			ml.currentBtn++
			if ml.currentBtn >= len(ml.buttons) {
				ml.currentBtn = len(ml.buttons)
			}
		case tl.KeyArrowUp:
			ml.currentBtn--
			if ml.currentBtn == 0 {
				ml.currentBtn = 1
			}
		}

		// highlight button
		ml.buttons[ml.currentBtn-1].SetColor(500)
		ml.buttonsText[ml.currentBtn-1].SetColor(tl.ColorBlack, 500)
	}
}

func addButton(l *tl.BaseLevel, ml *MenuListener, text string, x int, y int, width int) {
	button := tl.NewRectangle(x, y, width, 3, tl.ColorBlack)
	buttonText := tl.NewText(x+width/2-len(text)/2, y+1, text, tl.ColorWhite, tl.ColorBlack)

	ml.buttons = append(ml.buttons, button)
	ml.buttonsText = append(ml.buttonsText, buttonText)

	l.AddEntity(button)
	l.AddEntity(buttonText)
}

// SetupMainMenuLevel sets up the main level and returns it
func SetupMainMenuLevel() *tl.BaseLevel {
	level := tl.NewBaseLevel(tl.Cell{Fg: tl.ColorBlack, Bg: tl.ColorBlack, Ch: ' '})

	// add listener
	ml := &MenuListener{tl.NewEntity(0, 0, 0, 0), make([]*tl.Rectangle, 0), make([]*tl.Text, 0), 0}
	level.AddEntity(ml)

	// add background
	level.AddEntity(tl.NewRectangle(1, 1, 57, 23, tl.ColorWhite))

	// add title
	titleEntity := tl.NewEntityFromCanvas(7, 5, tl.CanvasFromString(BigTitleText))
	level.AddEntity(titleEntity)

	// add credit
	level.AddEntity(tl.NewText(24, 11, "by Freddie", tl.ColorBlack, tl.ColorWhite))

	// add buttons
	addButton(level, ml, "Create Game", 7, 13, 44)
	addButton(level, ml, "Join Game", 7, 17, 44)

	return level
}
