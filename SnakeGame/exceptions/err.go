package exceptions

import (
	"fmt"
	"log"
)

type SnakeError struct {
	Message string
	Err     error
}

func (e *SnakeError) Error() string {
	panic(fmt.Sprintf("%s: %v", e.Message, e.Err))
}

func CheckErrors(err error, errMsg string) error {
	if err != nil {
		log.Fatal(errMsg)
		log.Fatal(err)
		return &SnakeError{
			Message: errMsg,
			Err:     err,
		}

	}
	return nil

}
