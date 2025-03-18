package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Morgahl/flood/pkg/board"
	"github.com/Morgahl/flood/pkg/color"
	"github.com/Morgahl/flood/pkg/file"
	"github.com/Morgahl/flood/pkg/graph"
)

func main() {
	userInput := bufio.NewScanner(os.Stdin)
	scnr, err := file.New("./internal/fixtures/floodtest")
	if err != nil {
		log.Fatalln(err)
	}

NEXT:
	for scnr.Scan() {
		start := time.Now()
		base := scnr.Board()
		g := graph.New(base)
		b := board.New(base)
		b.Normalize()
		g.Normalize()
		fmt.Println(g)
		fmt.Println(time.Since(start))
		for {
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

			fmt.Printf(`
Choose:
1-6: Flood
next, n: Next
exit, x, quit, q: Exit
> `)

			if !userInput.Scan() {
				return
			}

			line := userInput.Text()
			switch line {
			case "exit", "x", "quit", "q":
				return

			case "next", "n":
				break NEXT

			case "1", "2", "3", "4", "5", "6":
				v := line[0] - '0'
				g.Flood(v)
				b.Flood(v)
				continue

			default:
				fmt.Println("Invalid input")
			}
		}
	}
}
