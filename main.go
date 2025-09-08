package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d received: %d\n", id, job)
	}
}

func main() {
	var num int
	flag.IntVar(&num, "n", 10, "number of workers")
	flag.Parse()

	jobs := make(chan int)

	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < num; i++ {
		go worker(i, jobs, &wg)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		defer close(jobs)
		counter := 1
		for {
			select {
			case <-exit:
				fmt.Println("Received interrupt signal. Shutting down...")
				return
			case jobs <- counter:
				counter++
			}
		}
	}()

	wg.Wait()
}
