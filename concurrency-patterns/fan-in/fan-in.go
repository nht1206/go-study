package fanin

import "sync"

func Merge(inChans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	send := func(in <-chan interface{}, wg *sync.WaitGroup) {
		defer wg.Done()
		for data := range in {
			out <- data
		}
	}
	var wg sync.WaitGroup
	wg.Add(len(inChans))
	for _, ch := range inChans {
		go send(ch, &wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
