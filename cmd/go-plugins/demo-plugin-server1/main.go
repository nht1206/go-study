package main

import (
	"fmt"
	"net"
	"strings"

	demo "github.com/nht1206/go-study/go-plugins/demo"
)

type driver struct {
}

func (d driver) Handle(req demo.Request) demo.Response {
	if len(req.Data) <= 0 {
		return demo.Response{
			IsContinue: false,
			Err:        "data is required",
		}
	}
	if strings.Contains(req.Data, "for demo1") {
		return demo.Response{
			Result:     fmt.Sprintf("%q is processed by demo plugin 1", req.Data),
			IsContinue: false,
		}
	}

	return demo.Response{
		IsContinue: true,
	}
}

func main() {
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	h := demo.NewHandler(&driver{}, demo.WithDemoPath("/demo"))

	if err := h.Serve(l); err != nil {
		panic(err)
	}
}
