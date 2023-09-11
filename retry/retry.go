package retry

import "time"

var DefaultRetryIntervals = []int{100, 250, 350, 500, 700, 1000}

func getRetryInterval(attempt int) int {
	if attempt < 0 || attempt > len(DefaultRetryIntervals)-1 {
		return DefaultRetryIntervals[len(DefaultRetryIntervals)-1]
	}
	return DefaultRetryIntervals[attempt]
}

func Retry(op func(attempt int) error, maxRetries int) error {
	attempt := 1
retry:
	err := op(attempt)
	if err != nil && attempt < maxRetries {
		time.Sleep(time.Duration(getRetryInterval(attempt)))
		attempt++
		goto retry
	}

	return err
}
