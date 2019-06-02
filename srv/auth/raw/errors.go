package raw

import (
	"errors"
)

// errors
var (
	ErrRepoNotFound = errors.New("repo not found")
	ErrRepoError    = errors.New("repo error, please retry")

	ErrUserNotExist = errors.New("user not exist")
	ErrRoleNotExist = errors.New("role not exist")
	ErrPermNotExist = errors.New("perm not exist")
	ErrUserHasExist = errors.New("user has exist")
	ErrRoleHasExist = errors.New("role has exist")
	ErrPermHasExist = errors.New("perm has exist")
	ErrPassWrong    = errors.New("pass wrong")

	ErrInvalidJWT = errors.New("jwt invalid")
	ErrSignJWT    = errors.New("jwt sign error")

	ErrNotUpdate    = errors.New("can't update")
	ErrUnknowAction = errors.New("unknow action")
)
