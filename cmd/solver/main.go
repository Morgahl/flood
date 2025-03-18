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

	"github.com/Morgahl/flood/pkg/file"

	// board "github.com/Morgahl/flood/pkg/board"
	// "github.com/Morgahl/flood/pkg/graph"
	// graph "github.com/Morgahl/flood/pkg/graph/v2"
	// graph "github.com/Morgahl/flood/pkg/graph/v3"
	// graph "github.com/Morgahl/flood/pkg/graph/v4"
	// graph "github.com/Morgahl/flood/pkg/graph/v5"
	graph "github.com/Morgahl/flood/pkg/graph/v6"
)

const LIMIT = 100_000

func main() {
	scnr, err := file.New("./internal/fixtures/floodtest")
	if err != nil {
		log.Fatalln(err)
	}
	defer scnr.Close()
	take, leave := tickets(runtime.GOMAXPROCS(0))

	start := time.Now()
	solutions := make([][]uint8, LIMIT)
	var solutionsCount atomic.Int64
	var solutionsLenSum atomic.Int64
	wg := &sync.WaitGroup{}
	wg.Add(LIMIT + 2)
	go func() {
		defer wg.Done()
		for i := 0; scnr.Scan() && i < LIMIT; i++ {
			take()
			g := graph.New(scnr.Board())
			g.Normalize()
			go func(i int, g *graph.Graph) {
				defer wg.Done()
				ng := g.Solve()
				solution := ng.Solution()
				leave()
				solutions[i] = solution
				solutionsCount.Add(1)
				solutionsLenSum.Add(int64(len(solution)))
			}(i, g)
		}
	}()

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			solCount := solutionsCount.Load()
			solLenSum := solutionsLenSum.Load()
			now := time.Now()
			elapsed := now.Sub(start)
			var avergage, solveRate float64
			var projection uint64
			var remaining time.Duration
			if solCount != 0 {
				avergage = float64(solLenSum) / float64(solCount)
				projection = uint64(float64(solLenSum) / float64(solCount) * LIMIT)
				remaining = time.Duration(int64(LIMIT-solCount) * int64(elapsed) / int64(solCount))
				solveRate = float64(solCount) / elapsed.Seconds()
			} else {
				remaining = time.Duration(LIMIT) * elapsed
			}
			progress := float64(solCount) / float64(LIMIT) * 100
			fmt.Printf(
				"%0.5d %0.8d %0.3f %d %0.3f %0.3f%% %s %s\n",
				solCount,
				solLenSum,
				avergage,
				projection,
				solveRate,
				progress,
				elapsed,
				remaining,
			)
			if solCount == LIMIT {
				return
			}
		}
	}()

	wg.Wait()
	fmt.Println(solutionsCount.Load())
	fmt.Println(solutionsLenSum.Load())
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
