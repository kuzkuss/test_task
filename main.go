package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

const (
	workersNum = 5
	searchStr = "Go"
)

func sumResults(wg *sync.WaitGroup, results chan int) {
	defer wg.Done()

	total := 0

	for res := range results {
		total += res
	}

	fmt.Printf("Total: %d\n", total)
}

func main() {
	wg := &sync.WaitGroup{}

	// channel for numbers of occurences
	results := make(chan int)

	wg.Add(1)
	go sumResults(wg, results)

	scanner := bufio.NewScanner(os.Stdin)

	// channel to control number of run workers
	limitWorkersCh := make(chan struct{}, workersNum)

	wgPath := &sync.WaitGroup{}

	// scan input data
	for scanner.Scan() {
		text := scanner.Text()

		wgPath.Add(1)
		go countStringOccurences(text, wgPath, limitWorkersCh, results)
	}
	wgPath.Wait()
	close(results)
	close(limitWorkersCh)

	wg.Wait()
}

