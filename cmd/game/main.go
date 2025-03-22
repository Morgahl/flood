package main

import (
	"log"

	"github.com/Morgahl/flood/pkg/game"
)

func main() {
	// game, err := game.New("./internal/fixtures/floodtest")
	// game, err := game.New("./internal/fixtures/ones")
	// game, err := game.New("./internal/fixtures/test")
	// game, err := game.New("./internal/fixtures/test_layered")
	// game, err := game.New("./internal/fixtures/1")
	// game, err := game.New("./internal/fixtures/5")
	game, err := game.New("./internal/fixtures/857")
	if err != nil {
		log.Fatalln(err)
	}
	if err := game.Run(); err != nil {
		log.Fatalln(err)
	}
}
