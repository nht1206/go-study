package main

import (
	"errors"
	"log"

	"github.com/nht1206/go-study/retry"
)

func main() {
	retry.Retry(func(attempt int) error {
		log.Println("attempt", attempt)

		if attempt == 3 {
			return nil
		}

		return errors.New("error")
	}, 5)
}
