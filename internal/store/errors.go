package store

import "errors"

var (
	ErrTimeoutExceeded   = errors.New("operation timed out")
	ErrOperationCanceled = errors.New("operation was canceled")
	ErrNotFound          = errors.New("not found")
)
