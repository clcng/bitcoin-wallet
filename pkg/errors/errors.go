package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	pkgerr "github.com/pkg/errors"
)

var (
	DefaultErrorCode = 1000
	PrintStacktrace  = false

	ErrorMap = map[int64]string{
		1000: "General API Error",
		1001: "Invalid Input Parameters",

		2001: "Fail to init key manager",
		2002: "Fail to get path key",
	}
)

type CodedError interface {
	error
	Code() int
}

//wrapper of pkg error
type withCode struct {
	cause error
	code  int
}

func (w *withCode) Code() int     { return w.code }
func (w *withCode) Error() string { return w.cause.Error() }
func (w *withCode) Cause() error  { return w.cause }

func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.Error())
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

//json.Marshaler interface
func (w *withCode) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"code":    w.Code(),
		"message": w.cause.Error(),
	}
	if PrintStacktrace {
		stacktrace := fmt.Sprintf("%+v", w.cause)
		m["stacks"] = strings.Split(stacktrace, "\n")
	}
	return json.Marshal(m)
}

func New(msg string) error {
	return &withCode{
		code:  DefaultErrorCode,
		cause: pkgerr.New(msg),
	}
}

func Coded(code int, msg string) error {
	return &withCode{
		code:  code,
		cause: pkgerr.New(msg),
	}
}

func Codedf(code int, format string, args ...interface{}) error {
	return &withCode{
		code:  code,
		cause: pkgerr.Errorf(format, args...),
	}
}

func Wrap(err error, msg string) error {
	code := DefaultErrorCode
	if cerr, ok := err.(CodedError); ok {
		code = cerr.Code()
	}
	return &withCode{
		code:  code,
		cause: pkgerr.Wrap(err, msg),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	code := DefaultErrorCode
	if cerr, ok := err.(CodedError); ok {
		code = cerr.Code()
	}
	return &withCode{
		code:  code,
		cause: pkgerr.Wrapf(err, format, args...),
	}
}

func CodedWrap(err error, code int, msg string) error {
	return &withCode{
		code:  code,
		cause: pkgerr.Wrap(err, msg),
	}
}

func CodedWrapf(err error, code int, format string, args ...interface{}) error {
	return &withCode{
		code:  code,
		cause: pkgerr.Wrapf(err, format, args...),
	}
}
