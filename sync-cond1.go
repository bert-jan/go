package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	cond := sync.NewCond(&mu)

	wg.Add(2)

	go func() {
		fmt.Println("Worker 1: Started")
		defer wg.Done()

		cond.L.Lock()
		fmt.Println("Worker 1: Waiting for condition")
		cond.Wait() // wait until condition is signaled
		fmt.Println("Worker 1: Condition met, proceeding")
		cond.L.Unlock()

		fmt.Println("Worker 1: Done")
	}()

	go func() {
		fmt.Println("Worker 2: Started")
		defer wg.Done()

		time.Sleep(3 * time.Second) // Simulate task before signaling

		cond.L.Lock()
		fmt.Println("Worker 2: Signaling condition")
		cond.Signal() // Signal the waiting goroutine
		cond.L.Unlock()

		fmt.Println("Worker 2: Done")
	}()

	wg.Wait()
	fmt.Println("All workers are done")
}
