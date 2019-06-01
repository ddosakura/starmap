package common

import (
	"github.com/micro/go-micro/errors"
)

// Forbid Action
func (s *RESTful) Forbid(e error) *RESTful {
	s.final = true
	s.err = e
	return s
}

// Forbid Action
func (a *DoAction) Forbid(e error) Action {
	a.stop = true
	a.s.err = e
	return a
}

// Logical for Role & Permission
type Logical int

// Logicals
const (
	LogicalAND Logical = iota
	LogicalOR
)

func aInB(a string, b []string) bool {
	for _, s := range b {
		if a == s {
			return true
		}
	}
	return false
}

// AInB test
func AInB(a []string, b []string) bool {
	for _, s := range a {
		if !aInB(s, b) {
			return false
		}
	}
	return true
}

func checkRuleRP(rules []string, logical Logical, result []string) error {
	switch logical {
	case LogicalAND:
		if AInB(rules, result) {
			return nil
		}
	case LogicalOR:
		for _, r := range result {
			if aInB(r, rules) {
				return nil
			}
		}
	}
	return errors.Forbidden(SrvName, "no permission")
}

// Role Check
func (s *RESTful) Role(rules []string, logical Logical) *RESTful {
	result, err := s.GetRoles()
	if err != nil {
		return s.Forbid(CleanErrResponse(SrvName, err, errors.Forbidden))
	}
	if e := checkRuleRP(rules, logical, result); e != nil {
		return s.Forbid(e)
	}
	return s
}

// Permission Check
func (s *RESTful) Permission(rules []string, logical Logical) *RESTful {
	result, err := s.GetPermissions()
	if err != nil {
		return s.Forbid(CleanErrResponse(SrvName, err, errors.Forbidden))
	}
	if e := checkRuleRP(rules, logical, result); e != nil {
		return s.Forbid(e)
	}
	return s
}

// Role Check
func (a *DoAction) Role(rules []string, logical Logical) Action {
	result, err := a.s.GetRoles()
	if err != nil {
		return a.Forbid(CleanErrResponse(SrvName, err, errors.Forbidden))
	}
	if e := checkRuleRP(rules, logical, result); e != nil {
		return a.Forbid(e)
	}
	return a
}

// Permission Check
func (a *DoAction) Permission(rules []string, logical Logical) Action {
	result, err := a.s.GetPermissions()
	if err != nil {
		return a.Forbid(CleanErrResponse(SrvName, err, errors.Forbidden))
	}
	if e := checkRuleRP(rules, logical, result); e != nil {
		return a.Forbid(e)
	}
	return a
}
