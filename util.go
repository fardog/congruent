package congruent

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

// BasicAuth creates an authentication string suitable for use in a header
func BasicAuth(u, p string) string {
	s := fmt.Sprintf("%s:%s", u, p)
	encoded := base64.StdEncoding.EncodeToString([]byte(s))

	return fmt.Sprintf("Basic %s", encoded)
}

func mergeHttpHeader(dest *http.Header, src *http.Header) {
	for k, va := range *src {
		dest.Del(k)

		for _, v := range va {
			dest.Add(k, v)
		}
	}
}

func mergeHttpHeaders(dest *http.Header, headers ...*http.Header) {
	for _, header := range headers {
		if header != nil {
			mergeHttpHeader(dest, header)
		}
	}
}
