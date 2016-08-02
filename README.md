# congruent

A Go utility library for testing that responses from multiple servers are
equivalent. Useful for regression testing when re-implementing/modifying an
existing service.

**Note:** This library is very WIP; you should expect _anything and everything_
to change at any time.

## Install

```
$ go get github.com/fardog/congruent
```

## Example

For a thorough example, see the [mkwords example][] which tests the [mkwords][]
public API against a locally running service.

## License

[MIT](./LICENSE)

[mkwords example]: ./example/mkwords_test.go
[mkwords]: https://mkwords.fardog.io
