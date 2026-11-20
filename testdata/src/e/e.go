package e

import (
	"errors"
	"io"
)

func do() error { return nil }

func wrap(errp *error) {}

// Issue #19.
func issue19() error {
	defer func() {}()

	err := do()
	if err != nil {
		return nil // want `error is not nil \(line 16\) but it returns nil`
	}

	return nil
}

// Issue #17: wrap may modify err through the pointer, so it is not reported.
func issue17() (err error) {
	defer wrap(&err)
	err = do()
	if err != nil {
		return nil
	}
	return nil
}

// Issue #18: the error is inspected by errors.Is before returning nil.
func errorsIs() error {
	err := do()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func errorsAs() error {
	err := do()
	if err != nil {
		var pe interface{ Timeout() bool }
		if errors.As(err, &pe) {
			return nil
		}
		return err
	}
	return nil
}

func errorsIsOnly() error {
	err := do()
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
