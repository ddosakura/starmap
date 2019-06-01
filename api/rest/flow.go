package rest

import (
	"context"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// Flow for RESTful API
type Flow struct {
	// req data (should not change by Middleware)
	Ctx    context.Context
	Req    *api.Request
	Res    *api.Response
	Rest   RESTfulType
	Params map[string]*api.Pair

	// cache
	AuthUserClient auth.UserService
	Token          *auth.UserToken
	Roles          []string
	Permissions    []string

	// other
	ctxForMiddle context.Context
	f            *Flow
	final        bool
	err          error
}

// Action for Flow
func (s *Flow) Action(t RESTfulType) RESTful {
	if s.final || s.Rest&t == 0 {
		return &pass{s}
	}

	s.final = true
	return &Flow{
		Ctx:            s.Ctx,
		Req:            s.Req,
		Res:            s.Res,
		Rest:           s.Rest,
		Params:         s.Params,
		AuthUserClient: s.AuthUserClient,
		Token:          s.Token,
		Roles:          s.Roles,
		Permissions:    s.Permissions,

		ctxForMiddle: s.ctxForMiddle,
		f:            s,
		final:        false,
	}
}

// Chain for Flow
func (s *Flow) Chain(fn Middleware) RESTful {
	if s.final {
		return s
	}
	e := fn(s.ctxForMiddle, s)
	if e != nil {
		s.final = true
		s.err = e
	}
	return s
}

// Done Flow
func (s *Flow) Done() RESTful {
	// update cache
	s.f.AuthUserClient = s.AuthUserClient
	s.f.Token = s.Token
	s.f.Roles = s.Roles
	s.f.Permissions = s.Permissions

	if s.final {
		s.f.final = true
		s.f.err = s.err
	}
	return s.f
}

// Final Flow
func (s *Flow) Final() error {
	if !s.final {
		return errors.MethodNotAllowed(SrvName, "%v is not allowed", s.Req.Method)
	}
	return s.err
}
