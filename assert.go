package congruent

import (
	"fmt"
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
					"Status for %s was %d, expected %d",
					r[i].Request.URL,
					r[i].StatusCode,
					status)
			}
		}
	}

	return nil
}
