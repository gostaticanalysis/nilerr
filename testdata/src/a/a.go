package a

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
		// return nil
		return nil // OK
	}

	return nil
}
