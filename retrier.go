package retrier

import "errors"

type Retrier struct {
	maxAttempt int
	attempt    int
	done       bool
	retrieable Retrieable
	loggerFunc func(error)
}

type Retrieable interface {
	Exec() error
}

var (
	ErrRetriableNil      = errors.New("Retriable object must not be nil")
	ErrMinimumMaxAttempt = errors.New("Max attempt should be greater than zero")
)

// New contstructing new retrier object
// ret are retriable object that has Exec method, maxAttempt is the maximum trying attempt before it ends, loggerFunc is optional if you want to put logger of any error does occurs in the process
func New(ret Retrieable, maxAttempt int, loggerFunc func(error)) (*Retrier, error) {
	if maxAttempt < 1 {
		return nil, ErrMinimumMaxAttempt
	}

	if ret == nil {
		return nil, ErrRetriableNil
	}

	return &Retrier{
		maxAttempt: maxAttempt,
		attempt:    0,
		retrieable: ret,
		loggerFunc: loggerFunc,
	}, nil
}

func (r *Retrier) run() {
	for !r.isDone() {
		err := r.do()
		if err == nil {
			r.done = true
		} else {
			if r.loggerFunc != nil {
				r.loggerFunc(err)
			}
		}
	}
}

func (r *Retrier) Start() {
	go r.run()
}

func (r *Retrier) do() error {
	r.attempt++
	return r.retrieable.Exec()
}

func (r *Retrier) isDone() bool {
	return r.attempt >= r.maxAttempt || r.done
}
