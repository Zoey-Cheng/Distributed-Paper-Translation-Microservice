package errutil

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin" // gin框架
)

// 错误结构体
type Error struct {
	httpCode int    // HTTP状态码
	code     int    // 业务错误码
	message  string // 错误信息
}

// 实现error接口
func (err *Error) Error() string {
	return fmt.Sprintf("httpCode: %v code: %d errMessage: %v", err.httpCode, err.code, err.message)
}

// 包装错误
type ErrorWrapper struct {
	error *Error  // 主错误
	errs  []error // 其他错误
}

// 实现error接口
func (err *ErrorWrapper) Error() string {
	return fmt.Sprintf("err: %+v more: %v", err.error, err.errs)
}

// 封装错误
func NewError(err *Error, errs ...error) *ErrorWrapper {
	return &ErrorWrapper{
		error: err,
		errs:  errs,
	}
}

// 响应错误
func ResponseError(ctx *gin.Context, err error, errs ...error) {

	log.Printf("response error: %v errs: %+v", err, errs) // 日志记录错误

	if r, ok := err.(*Error); ok { // 自定义错误,直接响应

		ctx.JSON(r.httpCode, gin.H{
			"errcode": r.code,
			"errmsg":  r.message,
		})
		return
	}

	if r, ok := err.(*ErrorWrapper); ok { // 封装的自定义错误,响应主错误

		ctx.JSON(r.error.httpCode, gin.H{
			"errcode": r.error.code,
			"errmsg":  r.error.message,
		})
		return
	}

	// 未知错误,响应默认错误
	ctx.JSON(UnknownError.httpCode, gin.H{
		"errcode": UnknownError.code,
		"errmsg":  UnknownError.message,
	})
}
