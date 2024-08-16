package uerror

import (
	"fmt"
)

// BusinessError 用于表示业务错误
type BusinessError struct {
	BasicError
}

func (e *BusinessError) Code() int32     { return e.code }
func (e *BusinessError) Message() string { return e.message }

func (e *BusinessError) Error() string {
	return fmt.Sprintf("business error: code: %d, message: %s", e.code, e.message)
}

// NewBusinessError 创建一个新的业务错误
func NewBusinessError(code int32, message string) *BusinessError {
	return &BusinessError{
		BasicError{
			errorCategory: ErrBiz,
			code:          code,
			message:       message,
		},
	}
}

// IsBusinessError 判断错误是否为业务错误
func IsBizError(err error) bool {
	for err != nil {
		switch v := err.(type) {
		case *BasicError:
			return v.errorCategory == ErrBiz
		case *BusinessError:
			return v.errorCategory == ErrBiz
		case interface{ Unwrap() error }:
			err = v.Unwrap()
		default:
			return false
		}
	}
	return false
}
