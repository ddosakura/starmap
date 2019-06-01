package rest

import (
	"context"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// LoadAuthService Wrapper
func LoadAuthService(loader func(ctx context.Context) (auth.UserService, bool)) Middleware {
	return func(ctx context.Context, s *Flow) error {
		AuthUserClient, ok := loader(s.Ctx)
		if !ok {
			return errors.InternalServerError(SrvName, "auth client not found")
		}
		s.AuthUserClient = AuthUserClient
		return nil
	}
}

// JWTCheck Wrapper (load UserInfo in s.Token)
func JWTCheck() Middleware {
	return func(ctx context.Context, s *Flow) error {
		Token := GetJWT(s.Req)
		token, err := s.AuthUserClient.Check(s.Ctx, &auth.UserToken{
			Token: Token,
		})
		if err != nil {
			return CleanErrResponse(SrvName, err, errors.Forbidden)
		}
		s.Token = token
		s.FreshJWT(token.Token)
		return nil
	}
}

// PCC - ParamCheck Config
type PCC struct {
	Must     bool
	Multi    bool
	DefaultV []string
}

// PCCS - PCC Group
type PCCS map[string]*PCC

// PCC in common use
var (
	PccEmptyStr = &PCC{DefaultV: []string{""}}
)

// ParamCheck Wrapper
func ParamCheck(pccs PCCS) Middleware {
	return func(ctx context.Context, s *Flow) error {
		for k, pcc := range pccs {
			pair := s.Params[k]
			if pair == nil || len(pair.Values) == 0 {
				if pcc.Must || pcc.Multi {
					return errors.BadRequest(SrvName, "need param `%s`", k)
				}
				if s.Params == nil || len(s.Params) == 0 {
					s.Params = make(map[string]*api.Pair)
				}
				s.Params[k] = &api.Pair{
					Values: pcc.DefaultV,
				}
				continue
			}
			if pcc.Multi && len(pair.Values) < 1 {
				return errors.BadRequest(SrvName, "param `%s` need multi-value", k)
			}
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
	return func(ctx context.Context, s *Flow) error {
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
	return func(ctx context.Context, s *Flow) error {
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
