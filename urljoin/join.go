// Package urljoin provides methods for joining URLs from url parts.
//
// When joining a URL, you need to keep track of slashes, ensuring that your
// joined URLs maintain path separators without doubling-up of slashes.
//
// This package is ported from the Node.js "url-join" library:
// https://github.com/jfromaniello/url-join
package urljoin

import (
	"regexp"
	"strings"
)

func normalize(url string) string {
	protocol := regexp.MustCompile(`:\/`)
	consecutive := regexp.MustCompile(`([^:\s])\/+`)
	trailing := regexp.MustCompile(`\/(\?|&|#[^!])`)
	query := regexp.MustCompile(`(\?.+)\?`)

	url = protocol.ReplaceAllString(url, "://")
	url = consecutive.ReplaceAllString(url, "$1/")
	url = trailing.ReplaceAllString(url, "$1")
	url = query.ReplaceAllString(url, "$1&")

	return url
}

// Join concatenates a series of URL fragments into a properly separated URL
func Join(parts ...string) string {
	return normalize(strings.Join(parts, "/"))
}
