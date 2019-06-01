package rest

import (
	"context"

	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

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
