package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job := <-jobs:
			fmt.Printf("worker %d received job %d\n", id, job)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	var num int
	flag.IntVar(&num, "n", 10, "number of workers")
	flag.Parse()

	jobs := make(chan int)

	var wg sync.WaitGroup
	wg.Add(num)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < num; i++ {
		go worker(ctx, i, jobs, &wg)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		counter := 1
		for {
			select {
			case <-exit:
				fmt.Println("Received interrupt signal. Shutting down...")
				cancel()
				close(jobs)
				return
			case jobs <- counter:
				counter++
			}
		}
	}()

	wg.Wait()
}
