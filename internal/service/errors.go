package service

import (
	"errors"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrAlreadyFollowed  = errors.New("already followed")
	ErrCannotFollowSelf = errors.New("cannot follow yourself")
	ErrUnauthorized     = errors.New("unauthorized")
)
