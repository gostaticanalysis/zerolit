# zerolit

[![pkg.go.dev][gopkg-badge]][gopkg]

`zerolit` finds return zero values but they are not literal.

## Install

You can get `zerolit` by `go install` command (Go 1.16 and higher).

```bash
$ go install github.com/gostaticanalysis/zerolit@latest
```

## How to use

`zerolit` run with `go vet` as below when Go is 1.12 and higher.

```bash
$ go vet -vettool=$(which zerolit) ./...
```

## Analyze with golang.org/x/tools/go/analysis

You can get analyzers of zerolit from [zerolit.Analyzers](https://pkg.go.dev/github.com/gostaticanalysis/zerolit/#Analyzers).
And you can use them with [unitchecker](https://golang.org/x/tools/go/analysis/unitchecker).

<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/gostaticanalysis/zerolit
[gopkg-badge]: https://pkg.go.dev/badge/github.com/gostaticanalysis/zerolit?status.svg
