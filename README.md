# nilerr

[![godoc.org][godoc-badge]][godoc]

`nilerr` finds code which returns nil even though it checks that error is not nil.

```go
func f() error {
	err := do()
	if err != nil {
		return nil // miss
	}
}
```

<!-- links -->
[godoc]: https://godoc.org/github.com/gostaticanalysis/nilerr
[godoc-badge]: https://img.shields.io/badge/godoc-reference-4F73B3.svg?style=flat-square&label=%20godoc.org

