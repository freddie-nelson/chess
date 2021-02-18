package main

import (
	// tl "github.com/JoelOtter/termloop"
	"github.com/containerd/console"
)

func main() {
	current := console.Current()

	ws, err := current.Size()
	if err != nil {
		return
	}

	ws.Height = 100
	ws.Width = 100
	current.Resize(ws)

	// screen := game.Screen()
	// screen.SetFps(30)

	// game.Start()
}
