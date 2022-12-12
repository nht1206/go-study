package fanout

func Double(in <-chan int) <-chan int64 {
	out := make(chan int64)

	go func() {
		for data := range in {
			out <- int64(data * 2)
		}
		close(out)
	}()

	return out
}
