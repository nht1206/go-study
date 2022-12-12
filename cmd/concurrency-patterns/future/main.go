package main

import (
	"context"

	"github.com/nht1206/go-study/concurrency-patterns/future"
)

func main() {
	printer := future.NewPrinter()
	printFuture := printer.Print()

	doubler := future.NewDoubler(3)
	doubleFuture := doubler.Double(context.Background(), printFuture.Sink())

	in := doubleFuture.Sink()
	for i := 1; i <= 10; i++ {
		in <- i
	}
	close(in)

	if err := doubleFuture.Wait(); err != nil {
		panic(err)
	}

	if err := printFuture.Wait(); err != nil {
		panic(err)
	}
}
