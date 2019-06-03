package check

import (
	"net/url"
)

// BaseRuleBuilder for check preparation
type BaseRuleBuilder struct {
	// Request
	Params url.Values // 参数值可以重名，但是大部分时候默认只看第一个
	//Param  func(string) string

	// Result
	Err error
}

// Rules Checking
func (b *BaseRuleBuilder) Rules(m *Must, pcms ...PCM) Loader {
	if b.Err != nil {
		return &BaseLoader{
			Err: b.Err,
		}
	}
	if m == nil {
		m = NoMust
	}
	entity, err := m.Test(b.Params)
	if err != nil {
		return &BaseLoader{
			Err: err,
		}
	}
	for _, pcm := range pcms {
		if e := pcm(entity); e != nil {
			return &BaseLoader{
				Err: e,
			}
		}
	}
	l := new(BaseLoader)
	l.entity = entity
	return l
}

type logical int

const (
	lNone logical = iota
	lAnd
	lOr
)

// Must Check
type Must struct {
	multi bool
	ks    []string
	a     *Must
	b     *Must
	l     logical
}

// M ust
func M(ks ...string) *Must {
	return &Must{
		ks: ks,
		l:  lNone,
	}
}

// And logical
func (m *Must) And(o *Must) *Must {
	return &Must{
		a: m,
		b: o,
		l: lAnd,
	}
}

// Or logical
func (m *Must) Or(o *Must) *Must {
	return &Must{
		a: m,
		b: o,
		l: lOr,
	}

}

// Multi logical
func (m *Must) Multi() *Must {
	m.multi = true
	return m
}

func (m *Must) test(cs url.Values) error {
	switch m.l {
	case lNone:
		for _, k := range m.ks {
			if cs[k] == nil || len(cs[k]) == 0 {
				return NewErrParamWrong(k, "must need")
			}
		}
	case lAnd:
		e := m.a.test(cs)
		if e != nil {
			return e
		}
		e = m.b.test(cs)
		if e != nil {
			return e
		}
	case lOr:
		e := m.a.test(cs)
		if e == nil {
			return nil
		}
		e = m.b.test(cs)
		if e != nil {
			return e
		}
	}
	return nil
}

// Test Params
func (m *Must) Test(vs url.Values) (Model, error) {
	if m != nil {
		if vs == nil {
			return nil, NewErrParamWrong("", "no params!")
		}
		if e := m.test(vs); e != nil {
			return nil, e
		}
	}
	v := make(Model, len(vs))
	if m.multi {
		for k := range vs {
			v[k] = vs[k]
		}
	} else {
		for k := range vs {
			v[k] = vs.Get(k)
		}
	}
	return v, nil
}

// NoMust Check
var NoMust *Must
