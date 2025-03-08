package zevs_errors

import (
	"fmt"
)

type PublicError struct {
	Text string
}

func (e *PublicError) Error() string {
	return fmt.Sprintf("%s", e.Text)
}
