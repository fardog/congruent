package urljoin

import (
	"testing"
)

type joinCase struct {
	result string
	expect string
}
type joinCases []joinCase

func TestJoin(t *testing.T) {
	cases := joinCases{
		joinCase{Join("http://www.google.com/", "foo/bar", "?test=123"),
			"http://www.google.com/foo/bar?test=123"},

		joinCase{Join("http://www.google.com/", "foo/bar", "?test=123"),
			"http://www.google.com/foo/bar?test=123"},

		joinCase{Join("http://www.google.com", "#!", "foo/bar", "?test=123"),
			"http://www.google.com/#!/foo/bar?test=123"},

		joinCase{Join("http:", "www.google.com/", "foo/bar", "?test=123"),
			"http://www.google.com/foo/bar?test=123"},

		joinCase{Join("http://", "www.google.com/", "foo/bar", "?test=123"),
			"http://www.google.com/foo/bar?test=123"},

		joinCase{Join("http:", "www.google.com///", "foo/bar", "?test=123"),
			"http://www.google.com/foo/bar?test=123"},

		joinCase{Join("http:", "www.google.com///", "foo/bar", "?test=123", "#faaaaa"),
			"http://www.google.com/foo/bar?test=123#faaaaa"},

		joinCase{Join("//www.google.com", "foo/bar", "?test=123"),
			"//www.google.com/foo/bar?test=123"},

		joinCase{Join("http:", "www.google.com///", "foo/bar", "?test=123", "?key=456"),
			"http://www.google.com/foo/bar?test=123&key=456"},

		joinCase{Join("http:", "www.google.com///", "foo/bar", "?test=123", "?boom=value", "&key=456"),
			"http://www.google.com/foo/bar?test=123&boom=value&key=456"},
	}

	for _, c := range cases {
		if c.result != c.expect {
			t.Errorf("expected %v, got %v", c.expect, c.result)
		}
	}
}
