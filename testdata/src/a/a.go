package a

import (
	"context"
	"errors"
	"log"
	"testing"
)

func f() error {
	err := do()
	if err != nil {
		return nil // want "error is not nil \\(line 11\\) but it returns nil"
	}

	if err := do(); err != nil {
		return nil // want "error is not nil \\(line 16\\) but it returns nil"
	}

	if err := do(); err != nil {
		//lint:ignore nilerr reason
		return nil // OK
	}

	return nil
}

func g() error {
	err := do()
	if err == nil {
		return err // want "error is nil \\(line 29\\) but it returns error"
	}

	if err := do(); err == nil {
		return err // want "error is nil \\(line 34\\) but it returns error"
	}

	bytes, err := do2()
	if err == nil {
		_ = bytes
		return err // want "error is nil \\(line 38\\) but it returns error"
	}

	if err := do(); err == nil {
		return errors.New("another error") // OK
	}

	if err := do(); err != nil {
		return errors.New(err.Error()) // OK, error is wrapped
	}

	if err := do(); err != nil {
		CustomLoggingFunc(err) // OK
		return nil
	}

	if err := do(); err != nil {
		Logf(context.Background(), "error: %+v", err) // OK
		return nil
	}

	if err := do(); err != nil {
		LogTypedf(context.Background(), "error: %+v", err) // OK
		return nil
	}

	if err := do(); err != nil {
		LogSinglef(context.Background(), "error: %+v", err) // OK
		return nil
	}

	if err := do(); err != nil {
		NewLogger().CustomLoggingFunc(err) // OK
		return nil
	}

	if err := do(); err != nil {
		//lint:ignore nilerr reason
		return nil // OK
	}

	return nil
}

func h() {
	f0 := func() error {
		for {
			if err := do(); err != nil {
				break
			}
		}
		return nil // want "error is not nil \\(line 88\\) but it returns nil"
	}
	_ = f0

	f1 := func(t *testing.T) error {
		for {
			if err := do(); err != nil {
				t.Fatal(err)
			}
		}
		return nil
	}
	_ = f1
}

func i() (error, error) {
	if err := do(); err != nil {
		return nil, nil // want "error is not nil \\(line 108\\) but it returns nil"
	}

	if err := do(); err != nil {
		return nil, err
	}

	if err := do(); err != nil {
		return err, nil
	}

	if err := do(); err != nil {
		return err, err
	}

	return nil, nil
}

func j() (interface{}, error) {
	if err := do(); err != nil {
		return nil, nil // want "error is not nil \\(line 128\\) but it returns nil"
	}

	if err := do(); err != nil {
		return nil, err
	}

	if err := do(); err != nil {
		return err, nil // want "error is not nil \\(line 136\\) but it returns nil"
	}

	if err := do(); err != nil {
		return err, err
	}

	return nil, nil
}

func k()  {
	if err := do(); err != nil {
		return
	}

	if err := do(); err == nil {
		return
	}
}

func do() error {
	return nil
}

func do2() ([]byte, error) {
	return nil, nil
}

func CustomLoggingFunc(err error) {
	log.Printf("%+v", err)
}

func Logf(ctx context.Context, msg string, args ...interface{}) {
	log.Printf(msg, args...)
}

func LogTypedf(ctx context.Context, msg string, args ...error) {
	log.Printf(msg, args[0])
}

func LogSinglef(ctx context.Context, msg string, arg interface{}) {
	log.Printf(msg, arg)
}

type logger int

func NewLogger() *logger {
	l := logger(0)
	return &l
}

func (l *logger) CustomLoggingFunc(err error) {
	log.Printf("%+v", err)
}
