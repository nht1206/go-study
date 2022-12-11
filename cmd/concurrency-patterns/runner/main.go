package main

import (
	"log"
	"time"

	"github.com/nht1206/go-study/concurrency-patterns/runner"
)

const timeout = 3 * time.Second

type task2 struct {
}

func (d task2) Execute() error {
	log.Println("task 1 executed")
	time.Sleep(2 * time.Second)
	return nil
}

func main() {
	r := runner.New(timeout)
	err := r.AddTask("task1", runner.ExecutorFunc(func() error {
		log.Println("task1 executed")
		time.Sleep(2 * time.Second)
		return nil
	}))
	if err != nil {
		log.Println(err)
		return
	}
	err = r.AddTask("task2", &task2{})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Starting work")
	if err := r.Start(); err != nil {
		log.Println(err.Error())
	}
	log.Println("Process ended")
}
