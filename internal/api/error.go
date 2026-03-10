package api

import (
	"errors"
)

var(
	ErrNoURL      = errors.New("NO_URL")
	ErrInvalidURL = errors.New("INVALID_URL")
	ErrReqFail    = errors.New("REQ_FAILED")
	ErrAuth       = errors.New("UNAUTHORIZED")
)