package berr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInternalServer = Newx("internal server error")
)

// Wrap convert regular error to BErr
func Wrap(err error, msg string) error {
	return wrapOrNew(err, CodeInternal, msg, getStackFrame(3))
}

// like Wrap but with code
func WrapWithCode(code ErrCode, err error, msg string) error {
	return wrapOrNew(err, code, msg, getStackFrame(3))
}

// like Wrap but without msg
func Wrapx(err error) error {
	return wrapOrNew(err, CodeInternal, "", getStackFrame(3))
}

// New BErr with stacktrace
func New(msg string) error {
	return new(errors.New(msg), CodeInternal, msg, getStackFrame(3))
}

// New BErr without stacktrace
func Newx(msg string) error {
	return new(errors.New(msg), CodeInternal, msg, nil)
}

type Error struct {
	code  ErrCode
	msg   string
	err   error
	stack Stack
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) chain() error {
	var buf strings.Builder
	for i := len(e.stack) - 1; i >= 0; i-- {
		buf.WriteString(fmt.Sprintf("%s|%d: ", e.stack[i].ShortFuncName(), e.stack[i].Line))
	}
	buf.WriteString(e.err.Error())
	return errors.New(buf.String())
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Code  ErrCode  `json:"code"`
		Msg   string   `json:"message"`
		Err   string   `json:"error"`
		Stack []string `json:"stack"`
	}{
		Code:  e.code,
		Msg:   e.msg,
		Err:   e.err.Error(),
		Stack: e.stack.Compact(),
	})
}

func wrapOrNew(err error, code ErrCode, msg string, sf *StackFrame) error {
	if err == nil {
		return nil
	}
	var e *Error
	switch {
	case errors.As(err, &e):
		return wrap(e, code, msg, sf)
	default:
		return new(err, code, msg, sf)
	}
}

// create new berr.Error
func new(err error, code ErrCode, msg string, sf *StackFrame) *Error {
	stack := make(Stack, 0)
	if sf != nil {
		stack = append(stack, *sf)
	}
	return &Error{
		code:  code,
		msg:   msg,
		err:   err,
		stack: stack,
	}
}

func wrap(err *Error, code ErrCode, msg string, sf *StackFrame) error {
	// chaining message
	if err.msg != "" && msg != "" {
		err.msg = msg + ": " + err.msg
	} else if err.msg == "" {
		err.msg = msg
	}

	// set code if undefined
	if err.code == CodeInternal && code != CodeInternal {
		err.code = code
	}

	// append stack
	if sf != nil {
		err.stack = append(err.stack, *sf)
	}

	return err
}

func Is(err error, target error) bool {
	return errors.Is(GetErrRoot(err), GetErrRoot(target))
}

func GetErrRoot(err error) error {
	var e *Error
	switch {
	case errors.As(err, &e):
		return e.err
	default:
		return err
	}
}

func GetErrChain(err error) error {
	var e *Error
	switch {
	case errors.As(err, &e):
		return e.chain()
	default:
		return err
	}
}

func GetCode(err error) ErrCode {
	var e *Error
	switch {
	case errors.As(err, &e):
		return e.code
	default:
		return CodeInternal
	}
}

func GetMsg(err error) string {
	var e *Error
	switch {
	case errors.As(err, &e):
		return e.msg
	case err == nil:
		return ""
	default:
		return err.Error()
	}
}

func GetMsgRoot(err error) string {
	var e *Error
	switch {
	case errors.As(err, &e):
		token := strings.Split(e.msg, ": ")
		if len(token) > 0 {
			return token[len(token)-1]
		}
		return ""
	case err == nil:
		return ""
	default:
		return err.Error()
	}
}

func GetStack(err error) Stack {
	var e *Error
	switch {
	case errors.As(err, &e):
		return e.stack
	default:
		return nil
	}
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
