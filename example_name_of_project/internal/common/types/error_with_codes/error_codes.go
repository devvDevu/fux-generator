package error_with_codes

import "strconv"

type ErrorCode int

func (e ErrorCode) Int() int {
	return int(e)
}

func (e ErrorCode) String() string {
	return strconv.Itoa(e.Int())
}

// cast error
const (
	_ ErrorCode = iota + 0
	CodeFailedToCast
)

var (
	ErrorFailedToCast = NewError("failed to cast object", CodeFailedToCast)
)

// cfg
const (
	_ ErrorCode = iota + 99
	CodeFailedToFindConfig
	CodeFailedToReadConfig
)

var (
	ErrorFailedToFindConfig = NewError("failed to find config", CodeFailedToFindConfig)
	ErrorFailedToReadConfig = NewError("failed to read config", CodeFailedToReadConfig)
)
