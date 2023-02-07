package payxgo_util

import (
	"fmt"
)

var (
	TypeError = NewError(1001, "传入的参数类型错误")
)

type PayxgoApiRequestError struct {
	/*
	   错误码
	*/
	PayxgoCode int `json:"code,omitempty" `

	/*
		错误信息
	*/
	PayxgoMsg string `json:"msg,omitempty" `
}

func NewError(code int, msg string) *PayxgoApiRequestError {
	e := &PayxgoApiRequestError{
		PayxgoMsg:  msg,
		PayxgoCode: code,
	}
	return e
}

func (e *PayxgoApiRequestError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.PayxgoCode, e.PayxgoMsg)
}
