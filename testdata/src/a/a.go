package a

import (
	"context"
	"errors"
	"log"
)

func do() error {
	return nil
}

func f() error {
	err := do()
	if err != nil {
		return nil // want "error is not nil but it returns nil"
	}

	if err := do(); err != nil {
		return nil // want "error is not nil but it returns nil"
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
		return err // want "error is nil but it returns error"
	}

	if err := do(); err == nil {
		return err // want "error is nil but it returns error"
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
