package c

func do() error { return nil }

// captured variable in a closure.
func closure() func() error {
	var err error
	return func() error {
		if err != nil {
			return nil // want `error is not nil \(line 9\) but it returns nil`
		}
		return nil
	}
}

func deferNamed() (err error) {
	defer func() {}()
	if e := do(); e != nil {
		return nil // want `error is not nil \(line 18\) but it returns nil`
	}
	return nil
}

func deferUnnamed() error {
	defer func() {}()
	if e := do(); e != nil {
		return nil // want `error is not nil \(line 26\) but it returns nil`
	}
	return nil
}

func deferNamedMulti() (n int, err error) {
	defer func() {}()
	if e := do(); e != nil {
		return 0, nil // want `error is not nil \(line 34\) but it returns nil`
	}
	return 0, nil
}

func deferIgnored() error {
	defer func() {}()
	if e := do(); e != nil {
		//lint:ignore nilerr reason
		return nil // OK
	}
	return nil
}

func deferReturnsError() error {
	defer func() {}()
	if e := do(); e != nil {
		return e
	}
	return nil
}

func deferErrIsNil() error {
	defer func() {}()
	err := do()
	if err == nil {
		return err // want `error is nil \(line 59\) but it returns error`
	}
	return nil
}

// The deferred function may change the named result, so it is not reported.
func deferModifiesResult() (err error) {
	defer func() {
		if err == nil {
			err = do()
		}
	}()
	if e := do(); e != nil {
		return nil
	}
	return nil
}

// named result returned by bare return is fine.
func namedBare() (err error) {
	if err = do(); err != nil {
		return
	}
	return
}

func generic[E error](e E) error {
	if err := do(); err != nil {
		return nil // want `error is not nil \(line 88\) but it returns nil`
	}
	return nil
}

func genericErr[E error](e E) error {
	var err error = e
	if err != nil {
		return nil // want `error is not nil but it returns nil`
	}
	return nil
}
