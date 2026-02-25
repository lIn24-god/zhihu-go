package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"` // 业务码：0 表示成功，非 0 表示具体错误
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据（成功时返回，失败时为 nil 或具体错误详情）
}

// Success 成功响应（HTTP 200，业务码 0）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Error 通用错误响应（可指定 HTTP 状态码和错误信息，业务码固定为 -1）
func Error(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, Response{
		Code: -1,
		Msg:  msg,
		Data: nil,
	})
}

// ErrorWithCode 带业务码的错误响应（更精细的错误分类）
func ErrorWithCode(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
