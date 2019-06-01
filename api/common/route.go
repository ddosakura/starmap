package common

import (
	"context"
	"encoding/json"

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
	Roles          []string
	Permissions    []string

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
		return s.Forbid(errors.InternalServerError(SrvName, "auth client not found"))
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
		return s.Forbid(CleanErrResponse(SrvName, err, errors.Forbidden))
	}
	s.Token = token
	s.FreshJWT(token.Token)
	return s
}

// GetRoles with cache
func (s *RESTful) GetRoles() ([]string, error) {
	if s.Roles == nil {
		result, err := s.AuthUserClient.Roles(s.Ctx, &auth.None{})
		if err != nil {
			return nil, err
		}
		s.Roles = result.Data
	}
	return s.Roles, nil
}

// GetPermissions with cache
func (s *RESTful) GetPermissions() ([]string, error) {
	//result, err := s.AuthUserClient.Permissions(s.Ctx, &auth.None{})
	//return result.Data, err

	// TODO: testing
	switch s.Rest {
	case POST:
		return []string{"user:insert"}, nil
	case DELETE:
		return []string{"user:select"}, nil
	case GET:
		return []string{"user:update", "user:select"}, nil
	case PUT:
		return []string{"user:delete", "user:insert"}, nil
	}
	return []string{}, nil
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
		return errors.MethodNotAllowed(SrvName, "%v is not allowed", s.Req.Method)
	}
	return CleanErrResponse(SrvName, s.err, errors.BadRequest)
}

// Action for RESTful
type Action interface {
	Role(rules []string, logical Logical) Action
	Permission(rules []string, logical Logical) Action
	Check(k string, multi bool, defautlV []string) Action // CheckParams
	Do(fn func(*RESTful) (interface{}, error)) *RESTful
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
			return a.Forbid(errors.BadRequest(SrvName, "need param `%s`", k))
		}
		if a.s.Params == nil || len(a.s.Params) == 0 {
			a.s.Params = make(map[string]*api.Pair)
		}
		a.s.Params[k] = &api.Pair{
			Values: defautlV,
		}
		return a
	}
	if multi && len(pair.Values) < 1 {
		return a.Forbid(errors.BadRequest(SrvName, "param `%s` need multi-value", k))
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
