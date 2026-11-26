package car

import "errors"

var (
	ErrCarNotFound   = errors.New("car not found")
	ErrCmdRepository = errors.New("unable to command repository")
)
