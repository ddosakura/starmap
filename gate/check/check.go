package check

import (
	"errors"
	"fmt"
)

// for init
var (
	Build BuildFunc

	// code 4xx
	ErrNotRESTfulRequest = errors.New("method not in POST/DELETE/GET/PUT")
	// ErrParamNotFound     = errors.New("param not found")

	// code 5xx
	ErrBuildConfig   = errors.New("config for BuildFunc is error")
	ErrRuleWrong     = errors.New("param check rule wrong")
	ErrAuthLoadCrash = errors.New("auto load crash")
)

// ErrParamWrong in check
//   code 4xx
type ErrParamWrong struct {
	name string
	msg  string
}

// NewErrParamWrong Builder
func NewErrParamWrong(name, msg string) ErrParamWrong {
	return ErrParamWrong{
		name,
		msg,
	}
}

func (e ErrParamWrong) Error() string {
	return fmt.Sprintf("param `%s` error: %s", e.name, e.msg)
}

// BuildFunc for init
type BuildFunc func(interface{}) RuleBuilder

// RuleBuilder for check preparation
type RuleBuilder interface {
	Rules(*Must, ...PCM) Loader
}

// Loader for params
type Loader interface {
	Load(...interface{}) error
}
