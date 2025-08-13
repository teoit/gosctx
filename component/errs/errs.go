package errs

import "errors"

var (
	ErrDataNotFound          = errors.New("data not found")
	ErrKeyCacheRedisNotEmpty = errors.New("key not empty")
)
