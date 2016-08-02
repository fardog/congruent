package congruent

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Responses []*Response

func (r Responses) StatusSame() error {
	if len(r) < 1 {
		return nil
	}

	return r.StatusEqual(r[0].StatusCode)
}

func (r Responses) StatusEqual(status int) error {
	for i, _ := range r {
		if i > 0 {
			if r[i].StatusCode != status {
				return fmt.Errorf(
					"(%s)%s: Status was %d, expected %d",
					r[i].Request.Method,
					r[i].Request.URL,
					r[i].StatusCode,
					status)
			}
		}
	}

	return nil
}

func (r Responses) HeaderSame() error {
	if len(r) < 1 {
		return nil
	}

	for _, resp := range r {
		headers := *resp.Headers
		for k, v := range headers {
			if err := r.HeaderEqual(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Responses) HeaderEqual(k string, v interface{}) error {
	switch v.(type) {
	case []string:
		return r.headerEqualWithArrayValue(k, v.([]string))
	case string:
		return r.headerEqualWithStringValue(k, v.(string))
	default:
		return fmt.Errorf("Expected string or string array as value")
	}
}

func (r Responses) headerEqualWithStringValue(k, v string) error {
	return r.headerEqualWithArrayValue(k, []string{v})
}

func (r Responses) headerEqualWithArrayValue(k string, v []string) error {
	k = http.CanonicalHeaderKey(k)

	for _, resp := range r {
		if resp == nil || resp.Request == nil {
			return fmt.Errorf(
				"Failed to check; a response was missing a request. This typically " +
					"happens when you fail to check an upstream error before an assertion")
		}
		url := resp.Request.URL
		method := resp.Request.Method

		header := *resp.Headers
		val, ok := header[k]
		if !ok {
			return fmt.Errorf(
				"(%s)%s: Expected header %v to have value %v, was nil",
				method, url, k, v)
		}

		if lval, lv := len(val), len(v); lval != lv {
			return fmt.Errorf(
				"(%s)%s: Expected header %v to have length %v, was %v",
				method, url, k, lv, lval)
		}

		for i, hv := range v {
			if hv != val[i] {
				return fmt.Errorf(
					"(%s)%s: Expected header %v to contain value %v at index %v, was %v",
					method, url, k, hv, i, val[i])
			}
		}
	}

	return nil
}

func (r Responses) BodySame() error {
	if len(r) < 2 {
		return nil
	}

	for i, resp := range r[1:] {
		if !bytesEqual(resp.Body, r[i].Body) {
			return fmt.Errorf(
				"(%s)%s:\nExpected body:\n  %s\nReceived body: \n  %s",
				resp.Request.Method, resp.Request.URL, cutBody(r[i].Body), cutBody(resp.Body))
		}
	}

	return nil
}

func bytesEqual(b, o []byte) bool {
	if len(b) != len(o) {
		return false
	}

	for i, c := range b {
		if c != o[i] {
			return false
		}
	}

	return true
}

func cutBody(b []byte) []byte {
	lstr := os.Getenv("CONGRUENT_MAX_DIFF")
	l, err := strconv.Atoi(lstr)
	if err != nil {
		l = 76
	}

	if len(b) > l {
		nb := b[:l]
		nb = append(nb, '.', '.', '.')

		return nb
	}

	return b
}
