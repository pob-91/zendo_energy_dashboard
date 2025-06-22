package errors

import "fmt"

type HttpError struct {
	StatusCode int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Http error with code: %d\n", e.StatusCode)
}
