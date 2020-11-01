package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/curlymon/flood/board"
	"github.com/curlymon/flood/scanner"
)

func main() {
	scnr, err := scanner.New("./floodtest")
	if err != nil {
		log.Fatalln(err)
	}
	defer scnr.Close()
	take, leave := tickets(runtime.GOMAXPROCS(0))

	start := time.Now()
	solutions := make([][]uint8, 100000)
	var solutionsCount int64
	wg := &sync.WaitGroup{}
	for i := 0; scnr.Scan(); i++ {
		wg.Add(1)
		take()
		go func(i int, g *board.Board) {
			defer wg.Done()
			defer leave()
			ng := g.Solve()
			solution := ng.Solution()
			solutions[i] = solution
			atomic.AddInt64(&solutionsCount, int64(len(solution)))
			fmt.Printf("%0.5d %0.7d %d\n", i, solutionsCount, solution)
		}(i, board.New(scnr.Board()))
	}

	wg.Wait()
	fmt.Println(solutionsCount)
	fmt.Println(time.Since(start))

	f, err := os.Create("solutions")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bf := bufio.NewWriterSize(f, 64*1024)
	defer bf.Flush()

	for _, solution := range solutions {
		for _, step := range solution {
			fmt.Fprintf(bf, "%d", step)
		}
		fmt.Fprintf(bf, "\n\n")
	}
}

func tickets(i int) (func(), func()) {
	ch := make(chan struct{}, i)
	for j := 0; j < i; j++ {
		ch <- struct{}{}
	}
	return func() {
			<-ch
		}, func() {
			ch <- struct{}{}
		}
}
