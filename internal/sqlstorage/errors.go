package sqlstorage

import (
	"fmt"
)

type AlreadyExistErr struct {
	key string
	Err error
}

func (ae *AlreadyExistErr) Error() string {
	return fmt.Sprintf("Already exist record with value: %v, %v ", ae.key, ae.Err)
}

func NewAlreadyExistErr(key string, err error) error {
	return &AlreadyExistErr{
		key: key,
		Err: err,
	}
}
