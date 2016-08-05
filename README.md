# congruent

A Go utility library for testing that responses from multiple servers are
equivalent. Useful for regression testing when re-implementing/modifying an
existing service.

See full docs at [godoc][].

**Note:** This library is very WIP; you should expect _anything and everything_
to change at any time. This library only targets testing APIs that return text
or JSON; content tests will not function for anything else.

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
[godoc]: https://godoc.org/github.com/fardog/congruent
