package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Morgahl/flood/pkg/board"
	"github.com/Morgahl/flood/pkg/color"
	"github.com/Morgahl/flood/pkg/file"

	// graph "github.com/Morgahl/flood/pkg/graph"
	// graph "github.com/Morgahl/flood/pkg/graph/v2"
	// graph "github.com/Morgahl/flood/pkg/graph/v3"
	// graph "github.com/Morgahl/flood/pkg/graph/v4"
	// graph "github.com/Morgahl/flood/pkg/graph/v5"
	// graph "github.com/Morgahl/flood/pkg/graph/v6"
	// graph "github.com/Morgahl/flood/pkg/graph/v7"
	// graph "github.com/Morgahl/flood/pkg/graph/v8"
	// graph "github.com/Morgahl/flood/pkg/graph/v9"
	graph "github.com/Morgahl/flood/pkg/graph/v10"
)

type Game struct {
	userInput     *bufio.Scanner
	puzzleScanner *file.Scanner
}

func New(puzzleFile string) (*Game, error) {
	scnr, err := file.New(puzzleFile)
	if err != nil {
		return nil, err
	}
	return &Game{
		userInput:     bufio.NewScanner(os.Stdin),
		puzzleScanner: scnr,
	}, nil
}

func (game *Game) Close() {
	if game.puzzleScanner != nil {
		game.puzzleScanner.Close()
	}
	if game.userInput != nil {
		game.userInput = nil
	}
}

func (game *Game) Run() error {
	defer game.Close()

	if err := game.selectLoop(); err != nil {
		return err
	}
	return nil
}

func (game *Game) input() string {
	if !game.userInput.Scan() {
		return ""
	}
	return game.userInput.Text()
}

func (game *Game) selectLoop() error {
	for {
		fmt.Printf(`
Choose:
next, n: Next
random, r: Random board
seeded, s: Seeded board
exit, x, quit, q: Exit
> `)

		switch game.input() {
		case "exit", "x", "quit", "q":
			return nil

		case "next", "n":
			if !game.puzzleScanner.Scan() {
				return nil
			}
			game.boardLoop(game.puzzleScanner.Board())
			continue

		case "random", "r":
			game.boardLoop(seededBoard(rand.Int63()))
			continue

		case "seeded", "s":
			fmt.Printf("Enter seed: ")
			pseed, err := strconv.ParseInt(game.input(), 10, 64)
			if err != nil {
				fmt.Printf("Invalid seed: %v\n", err)
				continue
			}
			game.boardLoop(seededBoard(pseed))
			continue

		default:
			fmt.Println("Invalid input")
		}
	}
}

func (game *Game) boardLoop(base [19][19]uint8) {
	start := time.Now()
	g := graph.New(base)
	b := board.New(base)
	b.Normalize()
	g.Normalize()
	initTook := time.Since(start)
	fmt.Println(initTook)
	for {
		fmt.Println(g)
		fmt.Println(b.ANSIString())
		fmt.Printf("Root Edge Values:\n")
		edgeValues := g.RootEdgeValues()
		for i := uint8(0); i < 7; i++ {
			count, exists := edgeValues[i]
			if !exists {
				continue
			}
			fmt.Printf("%s: %s\n", color.Colorize(i), count)
		}
	PROMPT:
		fmt.Printf(`
Choose:
1-6: Flood
solve, s: Solve
done, d
> `)

		line := game.input()
		switch line {
		case "1", "2", "3", "4", "5", "6":
			v := line[0] - '0'
			g.Flood(v)
			b.Flood(v)
			continue

		case "solve", "s":
			ng := g.Solve()
			for _, v := range ng.Solution() {
				b.Flood(v)
				g.Flood(v)
				fmt.Println(b.ANSIString())
			}
			fmt.Println(ng.Solution())
			goto PROMPT

		case "done", "d":
			return

		default:
			fmt.Println("Invalid input")
		}
	}
}

func seededBoard(seed int64) (board [19][19]uint8) {
	rnd := rand.New(rand.NewSource(seed))
	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			board[y][x] = uint8(rnd.Intn(6) + 1)
		}
	}
	return board
}
