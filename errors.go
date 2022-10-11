package main

import (
	"errors"
)

var (
	ErrEmptyCsv        = errors.New("csv file cannot be empty")
	ErrResponseEmpty   = errors.New("response empty")
	ErrInvalidMoveResp = errors.New("move response is not valid")
)
