package oaa

import (
	"net/http"
)

// Error 错误处理的结构体
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var (
	Success     = NewError(200, "处理完成！")
	ServerError = NewError(200, "系统异常，请稍后重试!")
	NotFound    = NewError(200, http.StatusText(http.StatusNotFound))
)

func OtherError(message string) *Error {
	return NewError(200, message)
}

func NewErrorCode(Code int) *Error {
	return NewError(Code, "")
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(Code int, msg string) *Error {
	return &Error{
		Code:    Code,
		Message: msg,
	}
}
