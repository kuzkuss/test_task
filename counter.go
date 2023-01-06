package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func countURL(path string) (int, error) {
	resp, err := http.Get(path)
	if err != nil {
		return 0, fmt.Errorf("error http get (URL %s): %v", path, err)
	}

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return 0, fmt.Errorf("status code is not successful: %d, path = %s", resp.StatusCode, path)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error read response body: %v, path = %s", err, path)
	}

	return bytes.Count(data, []byte(searchStr)), nil
}

func countFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("error read file %s: %v", path, err)
	}

	return bytes.Count((data), []byte(searchStr)), nil
}

func countStringOccurences(path string, wg *sync.WaitGroup, limitWorkersCh chan struct{}, results chan int) {
	limitWorkersCh <- struct{}{}

	defer wg.Done()
	defer func() { <-limitWorkersCh }()

	count := 0
	var err error

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		count, err = countURL(path)
	} else if path != "" {
		count, err = countFile(path)
	} else {
		return
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Count for %s: %d\n", path, count)
		results <- count
	}
}

