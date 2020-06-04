package a

import "errors"

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
		//lint:ignore nilerr reason
		return nil // OK
	}

	return nil
}
