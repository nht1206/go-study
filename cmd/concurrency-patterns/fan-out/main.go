package main

import (
	"log"
	"sync"

	fanout "github.com/nht1206/go-study/concurrency-patterns/fan-out"
)

func main() {
	in := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			in <- i
		}
		close(in)
	}()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		out1 := fanout.Double(in)
		for data := range out1 {
			log.Println(data)
		}
		wg.Done()
	}()

	go func() {
		out2 := fanout.Double(in)
		for data := range out2 {
			log.Println(data)
		}
		wg.Done()
	}()

	wg.Wait()
}
