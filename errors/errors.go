package errors

import "github.com/pkg/errors"

type customError struct {
	code          ErrorCode
	originalError error
	contextInfo   map[string]interface{}
}

func (e customError) Error() string {
	return e.originalError.Error()
}

// Cause 方法返回原始错误
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrap 方法用字符串封装错误
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Wrapf 方法用格式化字符串封装错误
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			code:          customErr.code,
			originalError: wrappedError,
			contextInfo:   customErr.contextInfo,
		}
	}
	return customError{code: ErrUnknown, originalError: wrappedError}
}

// AddErrorContext 方法为错误添加上下文
func AddErrorContext(err error, context map[string]interface{}) error {
	if customErr, ok := err.(customError); ok {
		return customError{code: customErr.code, originalError: customErr.originalError, contextInfo: context}
	}
	return customError{code: ErrUnknown, originalError: err, contextInfo: context}
}

// GetErrorContext 方法返回错误内容
func GetErrorContext(err error) map[string]interface{} {
	if customErr, ok := err.(customError); ok || customErr.contextInfo != nil {
		return customErr.contextInfo
	}
	return nil
}

// GetCode 方法返回错误类型
func GetCode(err error) ErrorCode {
	if customErr, ok := err.(customError); ok {
		return customErr.code
	}
	return ErrUnknown
}
