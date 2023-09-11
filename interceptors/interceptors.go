package interceptors

type InterceptorNextFunc[T any] func(t T)

type InterceptorFunc[T any] func(next InterceptorNextFunc[T]) InterceptorNextFunc[T]

func RunInterceptor[T any](data T, next InterceptorNextFunc[T], interceptors ...InterceptorFunc[T]) {
	for _, interceptor := range interceptors {
		next = interceptor(next)
	}

	next(data)
}
