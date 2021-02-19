package main

// Game stores values about the games current state
type Game struct {
	board        *[Size][Size]Spot
	ended        bool
	turn         int
	selectedSpot *Spot
	pickedSpot   *Spot
}

// SetSelectedSpot unselects the current spot and sets the new one
func (g *Game) SetSelectedSpot(file int, rank int) {
	// check if desired spot exists
	if file < 0 || file > Size-1 || rank < 0 || rank > Size-1 {
		return
	}

	// if a spot is already selected, unselect it
	if g.selectedSpot != nil {
		g.selectedSpot.selected = false
	}

	// set new spot
	g.selectedSpot = &g.board[file][rank]
	g.selectedSpot.selected = true
}

// ChangeSelectedSpot sets the new selected spot based on an offset from the current selected spot
func (g *Game) ChangeSelectedSpot(fileOff int, rankOff int) {
	current := g.selectedSpot
	if current == nil {
		return
	}

	file := current.file + fileOff
	rank := current.rank + rankOff

	g.SetSelectedSpot(file, rank)
}
