package main

import (
	"errors"
)
var (
	ErrResponseEmpty = errors.New("response empty")
	ErrInvalidMoveResp = errors.New("move response is not valid")
)