package handler

import (
	"context"

	"github.com/ddosakura/starmap/api/common"
	"github.com/ddosakura/starmap/api/rest"
	"github.com/ddosakura/starmap/api/sys/raw"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// Perm Handler
type Perm struct{}

// Entity API
func (*Perm) Entity(ctx context.Context, req *api.Request, res *api.Response) error {
	permService, ok := common.AuthPermFromContext(ctx)
	if !ok {
		return errors.InternalServerError(raw.SrvName, "auth client not found")
	}
	r := new(struct {
		Sub     string `json:"sub,omitempty"`
		Act     string `json:"act,omitempty"`
		Subject string `json:"subject,omitempty"`
		Action  string `json:"action,omitempty"`
		Detail  string `json:"detail,omitempty"`
	})

	return rest.REST(ctx, req, res).
		Chain(autoLoadAuthService).
		Chain(rest.JWTCheck()).
		Chain(rest.ParamCheck(rest.PCCS{
			"sub": rest.PccMust,
			"act": rest.PccMust,
		})).
		Chain(rest.ParamAutoLoad(nil, r)).
		// API
		// -> sub, act, detail
		// <- subject, action, detail
		Action(rest.POST).
		Chain(rest.ParamCheck(rest.PCCS{
			"detail": rest.PccMust,
		})).
		Chain(rest.PermCheck([]string{"perm:insert"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			perm, err := permService.Insert(ctx, &auth.PermInfo{
				Subject: r.Sub,
				Action:  r.Act,
				Detail:  r.Detail,
			})
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(perm)
		}).
		Done().
		// API
		// -> sub, act
		Action(rest.DELETE).
		Chain(rest.PermCheck([]string{"perm:delete"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			_, err := permService.Delete(ctx, &auth.PermInfo{
				Subject: r.Sub,
				Action:  r.Act,
			})
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(nil)
		}).
		Done().
		// API
		// -> sub, act
		// <- subject, action, detail
		Action(rest.GET).
		Chain(rest.PermCheck([]string{"perm:select"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			perm, err := permService.Select(ctx, &auth.PermInfo{
				Subject: r.Sub,
				Action:  r.Act,
			})
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(perm)
		}).
		Done().
		// API
		// -> sub, act, subject(opt), action(opt), detail(opt)
		// <- subject, action, detail
		Action(rest.PUT).
		Chain(rest.PermCheck([]string{"perm:update"}, rest.LogicalAND)).
		Chain(func(ctx context.Context, s *rest.Flow) error {
			perm, err := permService.Update(ctx, &auth.PermWrapper{
				Subject: r.Sub,
				Action:  r.Act,
				Perm: &auth.PermInfo{
					Subject: r.Subject,
					Action:  r.Action,
					Detail:  r.Detail,
				},
			})
			if err != nil {
				return rest.CleanErrResponse(raw.SrvName, err, errors.InternalServerError)
			}
			return s.Success(perm)
		}).
		Done().
		// Finish
		Final()
}
