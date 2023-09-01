package main

import (
	"bytes"
	"log"

	"github.com/nht1206/go-study/fqp"
)

func main() {
	scan("@request.auth.id != \"\" && (status = \"active\" || status = \"pending\")")
}

func scan(query string) {
	s := fqp.NewScanner(bytes.NewReader([]byte(query)))

	for s.HasNext() {
		t := s.Scan()
		if t.Type == fqp.TokenGroup {
			log.Println("group: ", t.Literal, "-----")
			scan(t.Literal)
			log.Println("endgroup------")
			continue
		}
		log.Println(t.Type, t.Literal)
	}

	if s.Err != nil {
		log.Println(s.Err)
	}
}
