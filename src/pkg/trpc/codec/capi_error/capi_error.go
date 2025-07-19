package capi_error

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type CapiError struct {
	ErrorCode  ErrorCode     `json:"ErrorCode"`
	ErrorMsg   string        `json:"Message"`
	Err        error         `json:"-"`
	Parameters []interface{} `json:"-"`
}

type CapiErrorInfo struct {
	ErrorCode ErrorCode
	ErrorMsg  []string
}

const LANGUAGE_CH = "chinese"
const LANGUAGE_EN = "english"

type ErrorCode string

var languageList = []string{LANGUAGE_CH, LANGUAGE_EN}

// NewError returns a wrapError instance
func NewError(code ErrorCode, msg string, err error, a ...interface{}) *CapiError {
	return &CapiError{
		ErrorCode:  code,
		ErrorMsg:   msg,
		Err:        err,
		Parameters: a,
	}
}

// NewError returns a wrapError instance
func NewErr(code ErrorCode, msg string, a ...interface{}) *CapiError {
	return &CapiError{
		ErrorCode:  code,
		ErrorMsg:   msg,
		Parameters: a,
	}
}

// ErrCode get the error code of WrapError
func ErrCode(err error) ErrorCode {
	if v, ok := err.(*CapiError); ok {
		return v.ErrorCode
	}
	return INTERNAL_ERROR_CODE
}

func (errorInfo *CapiError) Error() string {
	return errorInfo.ErrorMsg
}

func RegisterError(errorCode ErrorCode, msg []string) bool {
	if len(msg) != len(languageList) {
		return false
	}
	errorInfoMap[errorCode] = &CapiErrorInfo{
		ErrorCode: errorCode,
		ErrorMsg:  msg,
	}
	return true
}

// 获取错误信息
func (ec ErrorCode) ErrInfo(language string) *CapiError {
	index := slices.Index(languageList, language)
	if index == -1 {
		index = 0
	}
	if _, ok := errorInfoMap[ec]; !ok {
		return &CapiError{
			ErrorCode: INTERNAL_ERROR_CODE,
			ErrorMsg:  errorInfoMap[ec].ErrorMsg[index],
		}
	}
	return &CapiError{
		ErrorCode: ec,
		ErrorMsg:  errorInfoMap[ec].ErrorMsg[index],
	}
}

// Wrap returns a wrapped error with a stack trace at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg ...string) error {
	if err == nil {
		return nil
	}

	format := strings.Builder{}

	pc, file, line, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	format.WriteString(file[strings.LastIndex(file, "/")+1:] + "/")
	format.WriteString(name[strings.LastIndex(name, "/")+1:] + "():")
	format.WriteString(strconv.Itoa(line) + ": ")

	for _, t := range msg {
		format.WriteString(t + ": ")
	}

	if v, ok := err.(*CapiError); ok {
		v.ErrorMsg = fmt.Sprint(format.String(), v.Err.Error())
		format.WriteString("%w")
		v.Err = fmt.Errorf(format.String(), v.Err)
		return v
	}

	format.WriteString("%w")
	return fmt.Errorf(format.String(), err)
}

// Errorf formats according to a format specifier and returns an error with a stack trace at the point Errorf is called.
func Errorf(formats string, a ...interface{}) error {
	pc, file, line, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	err := fmt.Errorf(file[strings.LastIndex(file, "/")+1:]+"/"+name[strings.LastIndex(name, "/")+1:]+
		"():"+strconv.Itoa(line)+": "+formats, a...)

	return err
}
