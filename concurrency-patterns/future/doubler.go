package future

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type doubeFuture struct {
	in   chan<- int
	wait func() error
}

func (f doubeFuture) Sink() chan<- int {
	return f.in
}

func (f doubeFuture) Wait() error {
	return f.wait()
}

type Doubler struct {
	timeout time.Duration
	number  int
}

func NewDoubler(number int) *Doubler {
	return &Doubler{
		number: number,
	}
}

func (d Doubler) Double(ctx context.Context, out chan<- int64) *doubeFuture {
	in := make(chan int)

	double := func(data int) int64 {
		time.Sleep(2 * time.Second)
		return int64(data * 2)
	}

	doubleEg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < d.number; i++ {
		doubleEg.Go(func() (doublingErr error) {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case data, ok := <-in:
					if !ok {
						return nil
					}
					out <- double(data)
				}
			}
		})
	}

	return &doubeFuture{
		in: in,
		wait: func() error {
			if err := doubleEg.Wait(); err != nil {
				return err
			}
			close(out)
			return nil
		},
	}
}
