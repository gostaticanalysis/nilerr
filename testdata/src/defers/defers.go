package defers

func simpleFuncWithOneErr() error {
	if err := do(); err != nil {
		return err
	}
	if err := do(); err != nil {
		return nil // want "error is not nil \\(line 7\\) but it returns nil"
	}

	defer dont()

	if err := do(); err != nil {
		return nil // want "error is not nil \\(line 13\\) but it returns nil"
	}
	if err := do(); err != nil {
		return err
	}

	return nil
}

func simpleFuncWith2ndErr() ([]byte, error) {
	if err := do(); err != nil {
		return nil, err
	}
	if err := do(); err != nil {
		return nil, nil // want "error is not nil \\(line 27\\) but it returns nil"
	}

	defer dont()

	if err := do(); err != nil {
		return nil, nil // want "error is not nil \\(line 33\\) but it returns nil"
	}
	if err := do(); err != nil {
		return nil, err
	}

	return nil, nil
}

func simpleFuncWithoutErr() []byte {
	if err := do(); err != nil {
		return nil
	}
	if err := do(); err != nil {
		return nil
	}

	defer dont()

	if err := do(); err != nil {
		return nil
	}
	if err := do(); err != nil {
		return nil
	}

	return nil
}

func simpleFuncWithoutRet() {
	if err := do(); err != nil {
		return
	}
	if err := do(); err != nil {
		return
	}

	defer dont()

	if err := do(); err != nil {
		return
	}
	if err := do(); err != nil {
		return
	}

	return
}

// -----

func do() error {
	return nil
}

func dont() {}
