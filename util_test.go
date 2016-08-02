package congruent

import (
	"net/http"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	expected := "Basic dXNlcm5hbWU6cGFzc3dvcmQ="
	enc := BasicAuth("username", "password")
	if enc != expected {
		t.Errorf("expected %v, got %v", expected, enc)
	}
}

func TestMergeHTTPHeader(t *testing.T) {
	expected := http.Header{
		"Authorization":  []string{"stillnotvalidlol"},
		"Content-Type":   []string{"application/json"},
		"X-Custom-Thing": []string{"v custom"}}

	dest := http.Header{
		"Authorization": []string{"notvalidlol"},
		"Content-Type":  []string{"application/json"}}
	src := http.Header{
		"Authorization":  []string{"stillnotvalidlol"},
		"X-Custom-Thing": []string{"v custom"}}

	mergeHTTPHeader(&dest, &src)

	for k, va := range expected {
		h := dest.Get(k)
		if h != va[0] {
			t.Errorf("for header %v: got %v, expected %v", k, h, va[0])
		}
	}

	if ld, lv := len(dest), len(expected); ld != lv {
		t.Errorf("expected len did not match dest, %v != %v", ld, lv)
	}
}

func TestMergeHTTPHeaders(t *testing.T) {
	expected := http.Header{
		"Authorization":      []string{"stillnotrightlol"},
		"Content-Type":       []string{"application/json"},
		"X-Custom-Thing":     []string{"v custom"},
		"X-Custom-Thing-Two": []string{"so custom"}}

	dest := http.Header{
		"Authorization": []string{"notvalidlol"},
		"Content-Type":  []string{"application/json"}}
	src1 := http.Header{
		"Authorization":  []string{"stillnotvalidlol"},
		"X-Custom-Thing": []string{"v custom"}}
	src2 := http.Header{
		"Authorization":      []string{"stillnotrightlol"},
		"X-Custom-Thing-Two": []string{"so custom"}}

	mergeHTTPHeaders(&dest, &src1, &src2)

	for k, va := range expected {
		h := dest.Get(k)
		if h != va[0] {
			t.Errorf("for header %v: got %v, expected %v", k, h, va[0])
		}
	}

	if ld, lv := len(dest), len(expected); ld != lv {
		t.Errorf("expected len did not match dest, %v != %v", ld, lv)
	}
}
