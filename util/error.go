package util

import (
	"fmt"
	"time"
)

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() (s string) {
	s = fmt.Sprintf("at %v, %s", e.When, e.What)
	return
}
