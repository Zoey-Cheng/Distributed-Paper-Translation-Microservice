package errutil

import "net/http"

// 请求参数错误
var RequestParamError = &Error{
	httpCode: http.StatusBadRequest,
	code:     40000,
	message:  "请求参数错误",
}

// 未授权错误
var UnauthorizedError = &Error{
	httpCode: http.StatusUnauthorized,
	code:     40001,
	message:  "未授权,先登录",
}

// 授权失败错误
var GrantError = &Error{
	httpCode: http.StatusBadRequest,
	code:     40002,
	message:  "授权失败",
}

// 文件不存在错误
var FileNotExistError = &Error{
	httpCode: http.StatusBadRequest,
	code:     40003,
	message:  "文件不存在",
}

// 服务数据库错误
var ServerDBError = &Error{
	httpCode: http.StatusInternalServerError,
	code:     50000,
	message:  "服务数据库错误",
}

// 未知错误
var UnknownError = &Error{
	httpCode: http.StatusInternalServerError,
	code:     50001,
	message:  "未知错误",
}

// AI异常错误
var OpenAIError = &Error{
	httpCode: http.StatusInternalServerError,
	code:     50002,
	message:  "AI异常",
}
