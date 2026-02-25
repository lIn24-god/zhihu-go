package service

import "errors"

var (
	// 用户相关
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrCannotFollowSelf  = errors.New("cannot follow yourself")
	ErrAlreadyFollowed   = errors.New("already followed")

	// 文章相关
	ErrPostNotFound = errors.New("post not found")
	ErrPostNotOwned = errors.New("you do not own this post")

	// 频率限制
	ErrTooFrequent = errors.New("too frequent")

	// 其他
	ErrPermissionDenied = errors.New("permission denied") //权限不足

	//禁言相关
	ErrCannotMuteAdmin = errors.New("cannot mute an admin")
	ErrMuteFailed      = errors.New("failed to mute user")
	ErrUserMuted       = errors.New("user is muted") // 用户被禁言
)
