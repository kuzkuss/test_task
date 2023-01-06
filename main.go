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
	scanner := bufio.NewScanner(os.Stdin)

	// channel to control number of run workers
	limitWorkersCh := make(chan struct{}, workersNum)

	// channel for numbers of occurences
	results := make(chan int)

	wgRes := &sync.WaitGroup{}

	wgRes.Add(1)
	go sumResults(wgRes, results)

	wgStrs := &sync.WaitGroup{}

	// scan input data
	for scanner.Scan() {
		text := scanner.Text()

		wgStrs.Add(1)
		go countStringOccurences(text, wgStrs, limitWorkersCh, results)
	}
	wgStrs.Wait()
	close(results)
	close(limitWorkersCh)

	wgRes.Wait()
}

