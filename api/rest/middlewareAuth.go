package rest

import (
	"context"

	auth "github.com/ddosakura/starmap/srv/auth/proto"
	client "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Order: LoadAuthService -> JWTCheck -> RoleCheck/PermCheck/SuperRole

// AuthService needed to load
type AuthService interface {
	Check(ctx context.Context, in *auth.UserToken, opts ...client.CallOption) (*auth.UserToken, error)
	Roles(ctx context.Context, in *auth.Identity, opts ...client.CallOption) (*auth.Result, error)
	Perms(ctx context.Context, in *auth.Identity, opts ...client.CallOption) (*auth.Result, error)
}

// LoadAuthService Wrapper
func LoadAuthService(loader func(ctx context.Context) (AuthService, bool)) Middleware {
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

func roleLoader(ctx context.Context, s *Flow) error {
	if s.Roles == nil {
		result, err := s.AuthUserClient.Roles(s.Ctx, &auth.Identity{
			UUID: s.Token.User.UUID,
		})
		if err != nil {
			return CleanErrResponse(SrvName, err, errors.Forbidden)
		}
		s.Roles = result.Data
	}
	return nil
}

// RoleCheck Wrapper
func RoleCheck(rules []string, logical Logical) Middleware {
	return func(ctx context.Context, s *Flow) error {
		if err := roleLoader(ctx, s); err != nil {
			return err
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
			result, err := s.AuthUserClient.Perms(s.Ctx, &auth.Identity{
				UUID: s.Token.User.UUID,
			})
			if err != nil {
				return CleanErrResponse(SrvName, err, errors.Forbidden)
			}
			s.Perms = result.Data
		}

		if e := checkRuleRP(rules, logical, s.Perms); e != nil {
			return e
		}
		return nil
	}
}

// RoleLevelCheck Wrapper
func RoleLevelCheck(id string) Middleware {
	return func(ctx context.Context, s *Flow) error {
		if err := roleLoader(ctx, s); err != nil {
			return err
		}

		result, err := s.AuthUserClient.Roles(s.Ctx, &auth.Identity{
			UUID: id,
		})
		if err != nil {
			return CleanErrResponse(SrvName, err, errors.Forbidden)
		}

		// Lv.1 > Lv.9
		if RoleLevel(s.Roles) > RoleLevel(result.Data) {
			return errors.Forbidden(SrvName, "your role rating is too low")
		}
		return nil
	}
}
