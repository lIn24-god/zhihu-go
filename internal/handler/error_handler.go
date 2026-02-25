package handler

import (
	"errors"
	"net/http"

	"zhihu-go/internal/service"
	"zhihu-go/pkg/errcode"
	"zhihu-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// HandleError 将 service 层返回的错误转换为标准 API 错误响应
func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		response.ErrorWithCode(c, http.StatusNotFound, errcode.UserNotFound, "User not found")
	case errors.Is(err, service.ErrUserAlreadyExists):
		response.ErrorWithCode(c, http.StatusConflict, errcode.UserAlreadyExists, "User already exists")
	case errors.Is(err, service.ErrInvalidPassword):
		// 出于安全考虑，通常不区分“用户不存在”和“密码错误”，但为了演示这里分开
		response.ErrorWithCode(c, http.StatusUnauthorized, errcode.InvalidPassword, "Invalid password")
	case errors.Is(err, service.ErrUnauthorized):
		response.ErrorWithCode(c, http.StatusUnauthorized, errcode.Unauthorized, "Unauthorized")
	case errors.Is(err, service.ErrPostNotFound):
		response.ErrorWithCode(c, http.StatusNotFound, errcode.PostNotFound, "Post not found")
	case errors.Is(err, service.ErrPostNotOwned):
		response.ErrorWithCode(c, http.StatusForbidden, errcode.PostNotOwned, "You do not own this post")
	case errors.Is(err, service.ErrTooFrequent):
		response.ErrorWithCode(c, http.StatusTooManyRequests, errcode.TooFrequent, "Too many requests")
	case errors.Is(err, service.ErrPermissionDenied):
		response.ErrorWithCode(c, http.StatusForbidden, errcode.PermissionDenied, "Permission denied")
	case errors.Is(err, service.ErrCannotMuteAdmin):
		response.ErrorWithCode(c, http.StatusForbidden, errcode.CannotMuteAdmin, "Cannot mute an admin user")
	case errors.Is(err, service.ErrUserMuted):
		response.ErrorWithCode(c, http.StatusForbidden, errcode.UserMuted, "User is muted")
	default:
		// 未知错误（系统内部错误）
		response.ErrorWithCode(c, http.StatusInternalServerError, -1, "Internal server error")
		// 这里可以记录日志
		// log.Printf("unhandled error: %v", err)
	}
}
