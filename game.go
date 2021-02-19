package main

// Game stores values about the games current state
type Game struct {
	board *Board
	ended bool
	turn  int
}
