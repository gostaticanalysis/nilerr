package d

import (
	"net/url"
	"time"
)

// issue 16: importing standard packages such as time must not cause internal errors.
func parseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return d, nil
}

// issue 12: named result with an unrelated error variable.
func unescape(query string) (err error) {
	query, qerr := url.QueryUnescape(query)
	if qerr != nil {
		return nil // want `error is not nil \(line 19\) but it returns nil`
	}
	return nil
}

// issue 8: lint:ignore comments must suppress the diagnostic.
func ignoredLine() error {
	_, err := time.ParseDuration("1s")
	if err != nil {
		return nil //lint:ignore nilerr reason
	}
	return nil
}

func ignoredPrevLine() error {
	_, err := time.ParseDuration("1s")
	if err != nil {
		//lint:ignore nilerr reason
		return nil
	}
	return nil
}
