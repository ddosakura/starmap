package rest

import (
	"context"
	"encoding/json"

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
	PccMust     = &PCC{Must: true}
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

// ParamModify for ParamAutoLoad
type ParamModify map[string]func([]string) interface{}

// ParamAutoLoad Wrapper
func ParamAutoLoad(modify ParamModify, entity interface{}) Middleware {
	return func(ctx context.Context, s *Flow) error {
		ps := make(map[string]interface{})
		for k, p := range s.Params {
			if len(p.Values) > 0 {
				if modify == nil || modify[k] == nil {
					ps[k] = p.Values[0]
				} else {
					ps[k] = modify[k](p.Values)
				}
			}
		}
		bs, err := json.Marshal(ps)
		if err != nil {
			return err
		}
		return json.Unmarshal(bs, &entity)
	}
}
