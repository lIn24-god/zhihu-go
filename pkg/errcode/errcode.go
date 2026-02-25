package errcode

const (
	// 系统级错误码 (保留 0 表示成功)
	Success = 0

	// 用户模块 10000-10999
	UserNotFound      = 10001
	UserAlreadyExists = 10002
	InvalidPassword   = 10003
	Unauthorized      = 10004
	CannotMuteAdmin   = 10005
	UserMuted         = 10006

	// 文章模块 11000-11999
	PostNotFound = 11001
	PostNotOwned = 11002

	// 通用错误码
	TooFrequent      = 12001
	PermissionDenied = 12002
)
