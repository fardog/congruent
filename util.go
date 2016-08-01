package congruent

import (
	"encoding/base64"
	"fmt"
)

// BasicAuth creates an authentication string suitable for use in a header
func BasicAuth(u, p string) string {
	s := fmt.Sprintf("%s:%s", u, p)
	encoded := base64.StdEncoding.EncodeToString([]byte(s))

	return fmt.Sprintf("Basic %s", encoded)
}
