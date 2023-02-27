package main

import (
	"bytes"
	"log"

	"github.com/nht1206/go-study/fqp"
)

func main() {
	s := fqp.NewScanner(bytes.NewReader([]byte("(@collection.a == b && c == d) || a == c")))

	for s.HasNext() {
		t := s.Scan()

		log.Println(t.Type, t.Literal)
	}
}
