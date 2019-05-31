package common

import (
	"context"
	"encoding/json"

	"github.com/ddosakura/starmap/api/auth/raw"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
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
	Ctx    context.Context
	Req    *api.Request
	Res    *api.Response
	Rest   RESTfulType
	Params map[string]*api.Pair

	AuthUserClient auth.UserService
	Token          *auth.UserToken

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
		Ctx:  ctx,
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

// LoadAuthService for RESTful
func (s *RESTful) LoadAuthService(loader func(ctx context.Context) (auth.UserService, bool)) *RESTful {
	AuthUserClient, ok := loader(s.Ctx)
	if !ok {
		s.final = true
		s.err = errors.InternalServerError("starmap.api", "auth client not found")
	}
	s.AuthUserClient = AuthUserClient

	return s
}

// CheckJWT for RESTful
func (s *RESTful) CheckJWT() *RESTful {
	Token := GetJWT(s.Req)
	token, err := s.AuthUserClient.Check(s.Ctx, &auth.UserToken{
		Token: Token,
	})
	if err != nil {
		s.final = true
		s.err = CleanErrResponse(raw.SrvName, err, errors.BadRequest)
	}
	s.Token = token
	s.FreshJWT(token.Token)
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

// FreshJWT for RESTful
func (s *RESTful) FreshJWT(jwt string) {
	SetJWT(s.Res, jwt)
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
