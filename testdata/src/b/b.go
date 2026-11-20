package b

import "errors"

type result struct{ err error }

func do() error { return nil }

func sendToChannel(ch chan<- error) error {
	if err := do(); err != nil {
		ch <- err // error is passed on
		return nil
	}
	return nil
}

func storeToField(r *result) error {
	if err := do(); err != nil {
		r.err = err // error is saved
		return nil
	}
	return nil
}

func storeToPointer(p *error) error {
	if err := do(); err != nil {
		*p = err // error is saved
		return nil
	}
	return nil
}

func storeToSliceElem(errs []error) error {
	if err := do(); err != nil {
		errs[0] = err // error is saved
		return nil
	}
	return nil
}

func appendToSlice() ([]error, error) {
	var errs []error
	if err := do(); err != nil {
		errs = append(errs, err) // error is saved
		return errs, nil
	}
	return errs, nil
}

func goStatement() error {
	if err := do(); err != nil {
		go handle(err) // error is passed on
		return nil
	}
	return nil
}

func deferStatement() error {
	if err := do(); err != nil {
		defer handle(err) // error is passed on
		return nil
	}
	return nil
}

func closureCapture() error {
	if err := do(); err != nil {
		defer func() { handle(err) }() // error is captured
		return nil
	}
	return nil
}

func handle(err error) {}

func stillReported() error {
	if err := do(); err != nil {
		_ = errors.New("unrelated")
		return nil // want `error is not nil \(line 77\) but it returns nil`
	}
	return nil
}
