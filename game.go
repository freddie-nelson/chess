package main

// Game stores values about the games current state
type Game struct {
	board *[Size][Size]Spot
	ended bool
	turn  int
}
