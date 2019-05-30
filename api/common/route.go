package common

import (
	"context"
	"encoding/json"

	"github.com/ddosakura/starmap/api/auth/raw"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

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

// RESTful API
type RESTful struct {
	Req    *api.Request
	Res    *api.Response
	Rest   RESTfulType
	Params map[string]*api.Pair

	final bool
	err   error
}

// RESTfulType for API
type RESTfulType int

// RESTfulType(s)
const (
	POST   RESTfulType = 1 << iota // INSERT
	DELETE                         // DELETE
	GET                            // INSERT
	PUT                            // UPDATE
)

// dict
var (
	RESTfulTypeDict = map[string]RESTfulType{
		"GET":    GET,
		"get":    GET,
		"POST":   POST,
		"post":   POST,
		"delete": DELETE,
		"DELETE": DELETE,
		"put":    PUT,
		"PUT":    PUT,
	}
)

// REST Builder
func REST(ctx context.Context, req *api.Request, res *api.Response) *RESTful {
	s := &RESTful{
		Req:  req,
		Res:  res,
		Rest: RESTfulTypeDict[req.Method],

		final: false,
	}
	switch s.Rest {
	case GET, DELETE:
		s.Params = req.Get
	case POST, PUT:
		s.Params = req.Post
	}
	return s
}

// ACTION for RESTful
func (s *RESTful) ACTION(t RESTfulType) Action {
	if s.final || s.Rest&t == 0 {
		return &passAction{s}
	}

	s.final = true
	return &DoAction{s, false}
}

// Final of RESTful
func (s *RESTful) Final() error {
	if !s.final {
		return errors.MethodNotAllowed(raw.SrvName, "%v is not allowed", s.Req.Method)
	}
	return s.err
}

// Action for RESTful
type Action interface {
	Check(k string, multi bool, defautlV []string) Action
	Do(fn func(*RESTful) (interface{}, error)) *RESTful
}

type passAction struct {
	s *RESTful
}

func (a *passAction) Check(k string, multi bool, defautlV []string) Action {
	return a
}
func (a *passAction) Do(fn func(*RESTful) (interface{}, error)) *RESTful {
	return a.s
}

// DoAction for RESTful
type DoAction struct {
	s    *RESTful
	stop bool
}

// Check Params
func (a *DoAction) Check(k string, multi bool, defautlV []string) Action {
	if a.stop {
		return a
	}
	pair := a.s.Params[k]
	if pair == nil || len(pair.Values) == 0 {
		if defautlV == nil {
			a.stop = true
			a.s.err = errors.BadRequest(raw.SrvName, "need param `%s`", k)
		} else {
			a.s.Params[k].Values = defautlV
		}
		return a
	}
	if multi && len(pair.Values) < 1 {
		a.stop = true
		a.s.err = errors.BadRequest(raw.SrvName, "param `%s` need multi-value", k)
		return a
	}
	return a
}

// Do Action
func (a *DoAction) Do(fn func(*RESTful) (interface{}, error)) *RESTful {
	if !a.stop {
		if res, err := fn(a.s); err == nil {
			b, _ := json.Marshal(SUCCESS(res))
			a.s.Res.StatusCode = 200
			a.s.Res.Body = string(b)
		} else {
			a.s.err = err
		}
	}
	return a.s
}
