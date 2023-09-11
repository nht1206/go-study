package main

import (
	"log"

	"github.com/nht1206/go-study/interceptors"
)

func printString(s string) {
	log.Println(s)
}

func addMoreData(extraData string) interceptors.InterceptorFunc[string] {
	return func(next interceptors.InterceptorNextFunc[string]) interceptors.InterceptorNextFunc[string] {
		return func(s string) {
			s += extraData
			next(s)
		}
	}
}

func main() {
	interceptors.RunInterceptor[string](
		"\ntest",
		printString,
		addMoreData(" \nthis will be added next"),
		addMoreData(" \nthis will be added first"),
	)
}
