package runner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	ErrTimeout   = errors.New("the runnable time exceeded")
	ErrInterrupt = errors.New("received an interrupt signal")
)

// Runner runs a set of tasks within a given timeout and can be shut down
// on a operating system interrupt
type Runner struct {
	// interrupt channel reports a signal from the operating system
	interrupt chan os.Signal
	// completed channel reports that the processing is done
	completed chan error
	// timeout channel reports that the time has run out
	timeout <-chan time.Time
	// tasks holds a set of functions that are executed
	// synchronously in index order
	tasks map[string]TaskExecutor
}

// New returns a new ready-to-use runner
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal),
		completed: make(chan error),
		timeout:   time.After(d),
		tasks:     map[string]TaskExecutor{},
	}
}

func (r *Runner) AddTask(taskName string, executor TaskExecutor) error {
	_, existed := r.tasks[taskName]
	if existed {
		return fmt.Errorf("%q is added", taskName)
	}
	r.tasks[taskName] = executor

	return nil
}

func (r *Runner) Start() error {
	signal.Notify(r.interrupt, os.Interrupt)
	go func() {
		r.completed <- r.run()
	}()
	select {
	case err := <-r.completed:
		return err
	case <-r.timeout:
		return ErrTimeout
	}
}

func (r *Runner) run() error {
	for taskName, executor := range r.tasks {
		if r.gotInterrupt() {
			return ErrInterrupt
		}
		log.Printf("starting to execute %s\n", taskName)
		if err := executor.Execute(); err != nil {
			return err
		}
		log.Printf("%s ended\n", taskName)
	}
	return nil
}

func (r *Runner) gotInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}
