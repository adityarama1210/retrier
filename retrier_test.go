package retrier

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRetriable struct {
	cnt int
}

func (mr *mockRetriable) Exec() error {
	mr.cnt++

	if mr.cnt < 3 {
		return errors.New("test error")
	}

	return nil
}

func Test_New(t *testing.T) {
	type args struct {
		ret        Retrieable
		maxAttempt int
		loggerFunc func(err error)
	}
	tests := []struct {
		name            string
		args            args
		expectedRetrier Retrier
		expectedError   error
	}{
		{
			name: "Normal no error",
			args: args{
				ret:        new(mockRetriable),
				maxAttempt: 5,
			},
			expectedRetrier: Retrier{
				attempt:    0,
				maxAttempt: 5,
				retrieable: new(mockRetriable),
			},
			expectedError: nil,
		},
		{
			name: "Error max attempt less than 1",
			args: args{
				ret: new(mockRetriable),
			},
			expectedRetrier: Retrier{
				retrieable: new(mockRetriable),
			},
			expectedError: ErrMinimumMaxAttempt,
		},
		{
			name: "Error retriable nil",
			args: args{
				maxAttempt: 1,
			},
			expectedRetrier: Retrier{},
			expectedError:   ErrRetriableNil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			r, err := New(tc.args.ret, tc.args.maxAttempt, tc.args.loggerFunc)
			assert.Equal(tt, tc.expectedError, err)
			if r != nil {
				assert.Equal(tt, tc.expectedRetrier, *r)
			}
		})
	}
}

func Test_isDone(t *testing.T) {
	r := Retrier{
		attempt:    9,
		maxAttempt: 9,
	}

	assert.Equal(t, true, r.isDone())
}

func Test_do(t *testing.T) {
	r := Retrier{
		retrieable: new(mockRetriable),
	}

	r.do()
	assert.Equal(t, 1, r.attempt)
}

func Test_run(t *testing.T) {
	obj := new(mockRetriable)

	r := Retrier{
		maxAttempt: 5,
		retrieable: obj,
	}

	r.run()
	assert.Equal(t, 3, obj.cnt)
	assert.Equal(t, 3, r.attempt)
}
