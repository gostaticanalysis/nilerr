package a

import (
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
		CustomLoggingFunc(err) // OK, error is used (most probably, for logging)
		return nil
	}

	if err := do(); err != nil {
		NewLogger().CustomLoggingFunc(err) // OK, error is used (most probably, for logging)
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

type logger int

func NewLogger() *logger {
	l := logger(0)
	return &l
}

func (l *logger) CustomLoggingFunc(err error) {
	log.Printf("%+v", err)
}
