package main

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex
var totalTime time.Duration

func task(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()
	time.Sleep(time.Duration(id) * time.Second)

	elapsed := time.Since(startTime)

	mu.Lock()
	totalTime += elapsed
	mu.Unlock()

	fmt.Printf("Task %d took %v\n", id, elapsed)
}

func main() {
	var wg sync.WaitGroup

	numTasks := 5

	startTime := time.Now()

	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		go task(i, &wg)
	}

	wg.Wait()

	// Totale tijd die alle taken samen hebben geduurd
	totalElapsed := time.Since(startTime)

	fmt.Printf("\nTotale tijd voor alle taken: %v\n", totalElapsed)
	fmt.Printf("Totale verwerkte tijd (inclusief wachttijd per taak): %v\n", totalTime)
}
