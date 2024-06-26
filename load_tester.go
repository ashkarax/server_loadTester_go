package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	totalRequests   = 1000
	concurrentUsers = 100
	targetURL       = "http://localhost:8080"
)

func worker(id int, wg *sync.WaitGroup, ch chan struct{}, results chan<- time.Duration) {
	defer wg.Done()
	client := &http.Client{}

	for range ch {
		start := time.Now()
		resp, err := client.Get(targetURL)
		if err != nil {
			log.Printf("Worker %d: error: %v while requesting %s\n", id, err, targetURL)
			continue
		}
		log.Println(resp.StatusCode)
		resp.Body.Close()
		duration := time.Since(start)
		results <- duration
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan struct{}, totalRequests)
	results := make(chan time.Duration, totalRequests)

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go worker(i, &wg, ch, results)
	}

	start := time.Now()

	for i := 0; i < totalRequests; i++ {
		ch <- struct{}{}
	}

	close(ch)
	wg.Wait()
	close(results)

	elapsed := time.Since(start)

	var totalDuration time.Duration
	var maxDuration time.Duration
	var minDuration = time.Hour

	for result := range results {
		totalDuration += result
		if result > maxDuration {
			maxDuration = result
		}
		if result < minDuration {
			minDuration = result
		}
	}

	avgDuration := totalDuration / time.Duration(totalRequests)

	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Concurrent users: %d\n", concurrentUsers)
	fmt.Printf("Average response time: %v\n", avgDuration)
	fmt.Printf("Max response time: %v\n", maxDuration)
	fmt.Printf("Min response time: %v\n", minDuration)

}
