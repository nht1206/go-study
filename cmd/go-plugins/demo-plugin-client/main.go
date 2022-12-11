package main

import (
	"context"
	"log"
	"time"

	"github.com/nht1206/go-study/go-plugins/demo"
)

func main() {
	manager := demo.NewManager()

	demoPlugin1, err := demo.NewPluginCaller(&demo.PluginCallerConfig{
		Endpoint: "http://localhost:8081/demo",
		Timeout:  time.Duration(3 * time.Second),
	})
	if err != nil {
		panic(err)
	}

	manager.Register("demo1", demoPlugin1)

	demoPlugin2, err := demo.NewPluginCaller(&demo.PluginCallerConfig{
		Endpoint: "http://localhost:8082/demo",
		Timeout:  time.Duration(3 * time.Second),
	})
	if err != nil {
		panic(err)
	}

	manager.Register("demo2", demoPlugin2)

	data := []*demo.Request{
		{
			Data: "this is for demo1",
		},
		{
			Data: "this is for demo2",
		},
		{
			Data: "this is for nothing",
		},
	}

	plugins := make([]string, 0)
	plugins = append(plugins, "demo1", "demo2")

	for _, d := range data {
		res, errs := manager.CallPlugins(context.Background(), plugins, d)

		for plgName, err := range errs {
			log.Printf("calling plugin %q failed: %v\n", plgName, err)
		}

		log.Printf("%#v\n", res)
	}

}
