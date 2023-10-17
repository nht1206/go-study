package pool

import (
	"errors"
	"time"
)

var (
	ErrScheduleTimeout = errors.New("schedule: timed out")
)

type Pool struct {
	tasks chan func()
}

func NewPool(size, queueSize int) *Pool {
	pool := &Pool{
		tasks: make(chan func(), queueSize),
	}

	for i := 0; i < size; i++ {
		go pool.worker()
	}

	return pool
}

func (p *Pool) Schedule(task func()) {
	p.schedule(task, nil)
}

func (p *Pool) ScheduleTimeout(task func(), timeout time.Duration) error {
	return p.schedule(task, time.After(timeout))
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrScheduleTimeout
	case p.tasks <- task:
		return nil
	}
}

func (p *Pool) worker() {
	for task := range p.tasks {
		task()
	}
}
