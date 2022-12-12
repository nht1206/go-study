package main

import (
	"fmt"
	"log"

	fanin "github.com/nht1206/go-study/concurrency-patterns/fan-in"
)

func main() {
	ch1 := make(chan interface{})
	go func() {
		defer close(ch1)
		for i := 1; i <= 10; i++ {
			ch1 <- fmt.Sprintf("From ch1: %v", i)
		}
	}()

	ch2 := make(chan interface{})
	go func() {
		defer close(ch2)
		for i := 1; i <= 10; i++ {
			ch1 <- fmt.Sprintf("From ch2: %v", i)
		}
	}()

	ch3 := make(chan interface{})
	go func() {
		defer close(ch3)
		for i := 1; i <= 10; i++ {
			ch1 <- fmt.Sprintf("From ch3: %v", i)
		}
	}()

	mergedCh := fanin.Merge(ch1, ch2, ch3)

	for data := range mergedCh {
		log.Println("Received:", data)
	}

}
