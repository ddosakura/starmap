package rest

import (
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// --- Param Utils ---

// Pair2Str extract single string value
func Pair2Str(pair *api.Pair) (v string, ok bool) {
	if pair == nil {
		return "", false
	}
	if len(pair.Values) == 0 {
		return "", true
	}
	return pair.Values[0], true
}

// --- JWT Utils ---

// GetJWT from header
func GetJWT(r *api.Request) string {
	ps := r.Header["Token"]
	if ps == nil || len(ps.Values) == 0 {
		return ""
	}
	return ps.Values[0]
}

// SetJWT to header
func SetJWT(r *api.Response, jwt string) {
	if jwt == "" {
		return
	}
	if r.Header == nil {
		r.Header = make(map[string]*api.Pair, 1)
	}
	r.Header["Set-Token"] = new(api.Pair)
	r.Header["Set-Token"].Key = "Set-Token"
	r.Header["Set-Token"].Values = []string{jwt}
}

// --- Error Handle Utils ---

// CleanErrResponse to cancel circle Error
func CleanErrResponse(id string, e error, fn func(id, format string, a ...interface{}) error) error {
	if e == nil {
		return nil
	}
	E, ok := e.(*errors.Error)
	if ok {
		E.Id = id
		return E
	}
	return fn(id, e.Error())
}

// --- Rule Check Utils ---

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
