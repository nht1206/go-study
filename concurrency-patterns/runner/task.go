package runner

type TaskExecutor interface {
	Execute() error
}

type ExecutorFunc func() error

func (f ExecutorFunc) Execute() error {
	return f()
}
