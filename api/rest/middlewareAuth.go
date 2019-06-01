package rest

import (
	"context"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/micro/go-micro/errors"
	"github.com/kr/pretty"
)

// Order: LoadAuthService -> JWTCheck -> RoleCheck/PermCheck

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

// Logical for Role & Perm
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
			result, err := s.AuthUserClient.Roles(s.Ctx, &auth.None{
				UUID: s.Token.User.UUID,
			})
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

// PermCheck Wrapper
func PermCheck(rules []string, logical Logical) Middleware {
	return func(ctx context.Context, s *Flow) error {
		if s.Perms == nil {
			result, err := s.AuthUserClient.Perms(s.Ctx, &auth.None{
				UUID: s.Token.User.UUID,
			})
			if err != nil {
				return CleanErrResponse(SrvName, err, errors.Forbidden)
			}
			s.Perms = result.Data
		}

		pretty.Println(rules, logical, s.Perms)
		if e := checkRuleRP(rules, logical, s.Perms); e != nil {
			return e
		}
		return nil
	}
}