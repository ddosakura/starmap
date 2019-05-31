package common

import (
	api "github.com/micro/go-api/proto"
)

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
