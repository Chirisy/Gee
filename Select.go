package main

import (
	"sync"
	"time"
)

func test1(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(2 * time.Second)
	ch <- "test1"
}

func test2(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(5 * time.Second)
	ch <- "test2"
}

func main() {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go test1(ch1, wg)
	go test2(ch2, wg)

	select {
	case res := <-ch1:
		println("test1 returned", res)
	case res := <-ch2:
		println("test2 returned", res)
	}
	println("Waiting for goroutines to finish")
	wg.Wait()
	println("All goroutines finished")
	close(ch1)
	close(ch2)
}
