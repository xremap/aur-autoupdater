package internalerrors

import "errors"

var (
	ErrUnknownPackage error = errors.New("unknown package")
)
