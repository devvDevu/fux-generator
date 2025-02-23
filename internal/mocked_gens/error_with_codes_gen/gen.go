package error_with_codes_gen

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func GenerateErrorWithCodes(filePath string) error {
	err := os.WriteFile(filePath+"/error_codes.go", []byte(code1), 0644)
	if err != nil {
		logrus.Errorf("Error writing error with codes: %v", err)
		return err
	}

	cmd := exec.Command("goimports", "-w", filePath+"/error_codes.go")
	cmd.Run()

	logrus.Infof("generated code for %s", filePath+"/error_codes.go")

	err = os.WriteFile(filePath+"/error_with_codes.go", []byte(code2), 0644)
	if err != nil {
		logrus.Errorf("Error writing error with codes: %v", err)
		return err
	}

	cmd = exec.Command("goimports", "-w", filePath+"/error_with_codes.go")
	cmd.Run()

	logrus.Infof("generated code for %s", filePath+"/error_with_codes.go")

	return nil
}

const code1 = `
package common

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
)`

const code2 = `package common

import (
	"fmt"
)

type ErrorWithCodeI interface {
	String() string
	Error() string
	GetCode() ErrorCode
	GetMessage() string
	GetOperation() string
	Is(err error) bool
	SetOperation(operation string) *ErrorWithCode
}

type ErrorWithCode struct {
	code      ErrorCode
	message   string
	operation string
}

func NewError(message string, code ErrorCode) *ErrorWithCode {
	err := new(ErrorWithCode)
	err.code = code
	err.message = message

	return err
}

func NewErrorWithOperation(message string, code ErrorCode, operation string) *ErrorWithCode {
	err := new(ErrorWithCode)
	err.code = code
	err.message = message
	err.operation = operation

	return err
}

func (e *ErrorWithCode) SetOperation(operation string) *ErrorWithCode {
	return NewErrorWithOperation(e.message, e.code, operation)
}

func (e *ErrorWithCode) Error() string {
	return e.String()
}

func (e *ErrorWithCode) String() string {
	if e.operation != "" {
		return fmt.Sprintf("operation: %v, code: %v, message: %v", e.operation, e.code, e.message)
	}

	return fmt.Sprintf("code: %v, message: %v", e.code, e.message)
}

func (e *ErrorWithCode) GetCode() ErrorCode {
	return e.code
}

func (e *ErrorWithCode) GetMessage() string {
	return e.message
}

func (e *ErrorWithCode) GetOperation() string {
	return e.operation
}

func (e *ErrorWithCode) Is(err error) bool {
	errWith, errWithErr := ToErrorWithCode(err)
	if errWithErr != nil {
		return false
	}

	return e.GetCode() == errWith.code
}

func ToErrorWithCode(err error) (*ErrorWithCode, error) {
	if err == nil || !isErrorWithCode(err) {
		return nil, ErrorFailedToCast
	}

	return err.(*ErrorWithCode), nil
}

func isErrorWithCode(err error) bool {
	_, ok := err.(*ErrorWithCode)

	return ok
}`
