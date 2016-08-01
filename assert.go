package congruent

import (
	"fmt"
	"net/http"
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
					"%s: Status was %d, expected %d",
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

		header := *resp.Headers
		val, ok := header[k]
		if !ok {
			return fmt.Errorf(
				"%s: Expected header %v to have value %v, was nil", url, k, v)
		}

		if lval, lv := len(val), len(v); lval != lv {
			return fmt.Errorf(
				"%s: Expected header %v to have length %v, was %v", url, k, lv, lval)
		}

		for i, hv := range v {
			if hv != val[i] {
				return fmt.Errorf(
					"%s: Expected header %v to contain value %v at index %v, was %v",
					url, k, hv, i, val[i])
			}
		}
	}

	return nil
}
