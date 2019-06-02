package rest

import (
	"context"
	"encoding/json"

	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/errors"
)

// PCC - ParamCheck Config
type PCC struct {
	Must        bool
	Multi       bool
	DefaultV    []string
	Link        string
	LinkLogical Logical
	linkResult  map[string]bool
}

// PCCS - PCC Group
type PCCS map[string]*PCC

// PCC in common use
var (
	PccEmptyStr = &PCC{DefaultV: []string{""}}
	PccMust     = &PCC{Must: true}
)

// PccLabel Builder
func PccLabel(pcc *PCC, logical Logical) *PCC {
	return &PCC{
		Must:        pcc.Must,
		Multi:       pcc.Multi,
		DefaultV:    pcc.DefaultV,
		LinkLogical: logical,
		linkResult:  make(map[string]bool),
	}
}

// PccLink Builder
func PccLink(label string) *PCC {
	return &PCC{
		Link: label,
	}
}

func paramCheck(params map[string]*api.Pair, k string, pcc *PCC) error {
	pair := params[k]
	if pair == nil || len(pair.Values) == 0 {
		if pcc.Must || pcc.Multi {
			return errors.BadRequest(SrvName, "need param `%s`", k)
		}
		params[k] = &api.Pair{
			Values: pcc.DefaultV,
		}
		return nil
	}
	if pcc.Multi && len(pair.Values) < 1 {
		return errors.BadRequest(SrvName, "param `%s` need multi-value", k)
	}
	return nil
}

// ParamCheck Wrapper
func ParamCheck(pccs PCCS) Middleware {
	return func(ctx context.Context, s *Flow) error {
		if s.Params == nil || len(s.Params) == 0 {
			s.Params = make(map[string]*api.Pair)
		}
		for k, pcc := range pccs {
			if k[0] == '$' {
				continue
			}
			link := false
			if pcc.Link != "" {
				link = true
				pcc = pccs[pcc.Link]
			}

			if err := paramCheck(s.Params, k, pcc); err != nil {
				if link {
					pcc.linkResult[k] = false
					continue
				}
				return err
			}
			if link {
				pcc.linkResult[k] = true
			}
		}
	Label:
		for k, pcc := range pccs {
			if k[0] == '$' {
				switch pcc.LinkLogical {
				case LogicalAND:
					for k, r := range pcc.linkResult {
						if !r {
							return errors.BadRequest(SrvName, "param `%s` logical error", k)
						}
					}
					continue Label
				case LogicalOR:
					keys := make([]string, 0, len(pcc.linkResult))
					for k, r := range pcc.linkResult {
						if r {
							continue Label
						}
						keys = append(keys, k)
					}
					return errors.BadRequest(SrvName, "param `%v` logical error", keys)
				}
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
