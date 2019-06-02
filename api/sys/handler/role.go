package handler

import (
	"context"
	"fmt"

	// proto "github.com/ddosakura/starmap/api/sys/proto"
	"github.com/ddosakura/starmap/api/common"
	"github.com/ddosakura/starmap/api/rest"
	"github.com/ddosakura/starmap/api/sys/raw"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// Role Handler
type Role struct{}

// Entity API
func (*Role) Entity(ctx context.Context, req *api.Request, res *api.Response) error {
	// TODO: Role Entity API
	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		// API
		Action(rest.POST).
		Chain(rest.PermCheck([]string{"role:insert"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.DELETE).
		Chain(rest.PermCheck([]string{"role:delete"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.GET).
		Chain(rest.PermCheck([]string{"role:select"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// API
		Action(rest.PUT).
		Chain(rest.PermCheck([]string{"role:update"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			fmt.Println("M", s.Rest)
			return s.Success(fmt.Sprintf("M %v", s.Rest))
		}).
		Done().
		// Finish
		Final()
}

// Perm Modify API
func (*Role) Perm(ctx context.Context, req *api.Request, res *api.Response) error {
	roleService, ok := common.AuthRoleFromContext(ctx)
	if !ok {
		return errors.InternalServerError(raw.SrvName, "auth client not found")
	}

	m := new(auth.Modification)
	playload := func(ctx context.Context, s *rest.Flow) error {
		result, err := roleService.Perm(ctx, m)
		if err != nil {
			return err
		}
		return s.Success(result.Data)
	}

	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		Chain(rest.PermCheck([]string{"role:perm"}, rest.LogicalAND)).
		// API - Add perm for role
		// -> role, subject, action
		// <- []string
		Action(rest.POST).
		Chain(rest.ParamCheck(rest.PCCS{
			"role":    rest.PccRename(rest.PccMust, "name"),
			"subject": rest.PccMust,
			"action":  rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, m)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			m.Modify = auth.M_Add
			return nil
		}).
		Chain(playload).
		Done().
		// API - Del perm for role
		// -> role, subject, action
		// <- []string
		Action(rest.DELETE).
		Chain(rest.ParamCheck(rest.PCCS{
			"role":    rest.PccRename(rest.PccMust, "name"),
			"subject": rest.PccMust,
			"action":  rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, m)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			m.Modify = auth.M_Del
			return nil
		}).
		Chain(playload).
		Done().
		// API - All Perms
		// <- []string
		Action(rest.GET).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			m.Modify = auth.M_List
			return nil
		}).
		Chain(playload).
		Done().
		// Finish
		Final()
}
