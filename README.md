# nilerr

[![godoc.org][godoc-badge]][godoc]

`nilerr` finds codes which are returning nil and checking either error is not nil.

```go
func f() error {
	err := do()
	if err != nil {
		return nil // miss
	}
}
```

`nilerr` also finds codes which are returning error and checking either error is nil.

```go
func f() error {
	err := do()
	if err == nil {
		return err // miss
	}
}
```

`nilerr` ignores a code which has a miss with ignore comment.

```go
func f() error {
	err := do()
	if err != nil {
		//lint:ignore nilerr reason
		return nil // ignore
	}
}
```

<!-- links -->
[godoc]: https://godoc.org/github.com/gostaticanalysis/nilerr
[godoc-badge]: https://img.shields.io/badge/godoc-reference-4F73B3.svg?style=flat-square&label=%20godoc.org

