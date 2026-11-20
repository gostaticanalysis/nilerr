# nilerr

[![pkg.go.dev][gopkg-badge]][gopkg]

`nilerr` finds code which returns nil even though it checks that error is not nil.

```go
func f() error {
	err := do()
	if err != nil {
		return nil // miss
	}
}
```

`nilerr` also finds code which returns error even though it checks that error is nil.

```go
func f() error {
	err := do()
	if err == nil {
		return err // miss
	}
}
```

`nilerr` ignores code which has a miss with ignore comment.

```go
func f() error {
	err := do()
	if err != nil {
		//lint:ignore nilerr reason
		return nil // ignore
	}
}
```

## How to use
```
$ go install github.com/gostaticanalysis/nilerr/cmd/nilerr@latest
$ nilerr ./...
```

## Development

Development tools are managed by [aqua](https://aquaproj.github.io/) (versions and checksums are pinned in `aqua.yaml` / `aqua-checksums.json`) and tasks are run with [Task](https://taskfile.dev/).

```
$ aqua i       # install the tools
$ task --list  # list tasks
$ task test    # goimports + go test -race ./...
$ task ci      # run all checks executed on CI
```

GitHub Actions are pinned to commit SHAs with [pinact](https://github.com/suzuki-shunsuke/pinact) (`task pinact`).

### Release

Releases are managed by [tagpr](https://github.com/Songmu/tagpr). Merging the release PR creates a `vX.Y.Z` tag.
Pull requests and tags created with the default `GITHUB_TOKEN` do not trigger other workflows,
so tagpr uses a GitHub App installation token (`actions/create-github-app-token`).
Register the App ID as the `TAGPR_APP_ID` secret and its private key as the `TAGPR_APP_PRIVATE_KEY` secret.
The App needs Contents (write), Pull requests (write) and Metadata (read) permissions.

<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/gostaticanalysis/nilerr
[gopkg-badge]: https://pkg.go.dev/badge/github.com/gostaticanalysis/nilerr?status.svg
