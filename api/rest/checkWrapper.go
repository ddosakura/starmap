package rest

import (
	"context"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// LoadAuthService Wrapper
func LoadAuthService(loader func(ctx context.Context) (auth.UserService, bool)) Middleware {
	return func(s *Flow) error {
		AuthUserClient, ok := loader(s.Ctx)
		if !ok {
			return errors.InternalServerError(SrvName, "auth client not found")
		}
		s.AuthUserClient = AuthUserClient
		return nil
	}
}

// JWTCheck Wrapper
func JWTCheck() Middleware {
	return func(s *Flow) error {
		Token := GetJWT(s.Req)
		token, err := s.AuthUserClient.Check(s.Ctx, &auth.UserToken{
			Token: Token,
		})
		if err != nil {
			return CleanErrResponse(SrvName, err, errors.Forbidden)
		}
		s.Token = token
		SetJWT(s.Res, token.Token)
		return nil
	}
}

// ParamCheck Wrapper
func ParamCheck(k string, multi bool, defautlV []string) Middleware {
	return func(s *Flow) error {
		pair := s.Params[k]
		if pair == nil || len(pair.Values) == 0 {
			if defautlV == nil {
				return errors.BadRequest(SrvName, "need param `%s`", k)
			}
			if s.Params == nil || len(s.Params) == 0 {
				s.Params = make(map[string]*api.Pair)
			}
			s.Params[k] = &api.Pair{
				Values: defautlV,
			}
			return nil
		}
		if multi && len(pair.Values) < 1 {
			return errors.BadRequest(SrvName, "param `%s` need multi-value", k)
		}
		return nil
	}
}

// Logical for Role & Permission
type Logical int

// Logicals
const (
	LogicalAND Logical = iota
	LogicalOR
)

// RoleCheck Wrapper
func RoleCheck(rules []string, logical Logical) Middleware {
	return func(s *Flow) error {
		if s.Roles == nil {
			result, err := s.AuthUserClient.Roles(s.Ctx, &auth.None{})
			if err != nil {
				return CleanErrResponse(SrvName, err, errors.Forbidden)
			}
			s.Roles = result.Data
		}

		if e := checkRuleRP(rules, logical, s.Roles); e != nil {
			return e
		}
		return nil
	}
}

// PermissionCheck Wrapper
func PermissionCheck(rules []string, logical Logical) Middleware {
	return func(s *Flow) error {
		if s.Permissions == nil {
			result, err := s.AuthUserClient.Permissions(s.Ctx, &auth.None{})
			if err != nil {
				return CleanErrResponse(SrvName, err, errors.Forbidden)
			}
			s.Permissions = result.Data
		}

		if e := checkRuleRP(rules, logical, s.Permissions); e != nil {
			return e
		}
		return nil
	}
}
