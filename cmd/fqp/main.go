package main

import (
	"bytes"
	"log"

	"github.com/nht1206/go-study/fqp"
)

func main() {
	scan("(a == 1 && b == 2) || c == 3")
}

func scan(query string) {
	s := fqp.NewScanner(bytes.NewReader([]byte(query)))

	for s.HasNext() {
		t := s.Scan()
		if t.Type == fqp.TokenGroup {
			log.Println("group: ", t.Literal, "-----")
			scan(t.Literal[1 : len(t.Literal)-1])
			log.Println("endgroup------")
			continue
		}
		log.Println(t.Type, t.Literal)
	}

	if s.Err != nil {
		log.Println(s.Err)
	}
}
