package congruent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// DefaultDiffLength is the default number of bytes that will be returned in a
// diff. Right now this is only respected for the request Body diff, but may
// expand to others in the future.
const DefaultDiffLength = 76

// Responses is an array of Response pointers
type Responses []*Response

// StatusSame verifies that all responses have the same status codes; returns
// an error for the first mismatch, if not.
func (r Responses) StatusSame() error {
	if len(r) < 1 {
		return nil
	}

	return r.StatusEqual(r[0].StatusCode)
}

// StatusEqual verifies that all responses match a given status code; returns an
// error for the first mismatch, if not.
func (r Responses) StatusEqual(status int) error {
	for i := range r {
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

// HeaderSame verifies that all headers match for all responses; returns an
// error for the first mismatch, if not.
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

// HeaderEqual verifies that a single header of key `k` matches the value `v`;
// value is expected to be either a `string` or `[]string`.
func (r Responses) HeaderEqual(k string, v interface{}) error {
	switch v.(type) {
	case []string:
		return r.headerEqualWithArrayValue(k, v.([]string))
	case string:
		return r.headerEqualWithStringValue(k, v.(string))
	default:
		return fmt.Errorf("Expected `string` or `[]string` array as value")
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

// BodySame verifies that the response body was identical on all requests;
// returns an error for the first mismatch if not. This is a bytewise
// comparison.
// When an error occurs, the response bodies will be trimmed to
// `DefaultDiffLength`, but this can be overridden by setting the environment
// variable `CONGRUENT_MAX_DIFF` to an integer value.
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

// BodyContentSame ensures that response bodies are roughly equivalent JSON or
// strings; no other content types can be expected to be handled appropriately.
// In here, JSON is Unmarshal'd, Marshal'd, and then compared. This results in
// only the contents being taken into account, and things like newlines,
// indentation, and etc being ignored.
func (r Responses) BodyContentSame() error {
	var bodies [][]byte

	for i, resp := range r {
		var body interface{}
		if resp.Body == nil {
			bodies = append(bodies, []byte{})
			continue
		}

		if err := json.Unmarshal(resp.Body, &body); err != nil {
			// hacky: if unmarshalling fails, treat as a string
			fmt.Println(err)
			bodies = append(bodies, resp.Body)
			continue
		}

		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodies = append(bodies, b)

		if i > 0 && !bytesEqual(b, bodies[i-1]) {
			return fmt.Errorf(
				"(%s)%s:\nExpected body:\n  %s\nReceived body: \n  %s",
				r[i-1].Request.Method, r[i-1].Request.URL, cutBody(bodies[i-1]), cutBody(b))
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
		l = DefaultDiffLength
	}

	if len(b) > l {
		nb := b[:l]
		nb = append(nb, '.', '.', '.')

		return nb
	}

	return b
}
