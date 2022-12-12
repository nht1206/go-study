package future

import (
	"log"

	"golang.org/x/sync/errgroup"
)

type printFuture struct {
	in   chan<- int64
	wait func() error
}

func (f printFuture) Sink() chan<- int64 {
	return f.in
}

func (f printFuture) Wait() error {
	return f.wait()
}

type Printer struct {
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p Printer) Print() *printFuture {
	in := make(chan int64)

	var eg errgroup.Group

	eg.Go(func() error {
		for data := range in {
			log.Println(data)
		}
		return nil
	})

	return &printFuture{
		in: in,
		wait: func() error {
			return eg.Wait()
		},
	}
}
