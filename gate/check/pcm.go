package check

import "strconv"

// PCM - ParamCheck Middleware
type PCM func(Model) error

// Rename PCM
func Rename(o, n string) PCM {
	return func(m Model) error {
		v := m[o]
		delete(m, o)
		m[n] = v
		return nil
	}
}

// DefaultValue PCM
func DefaultValue(k, v string) PCM {
	return func(m Model) error {
		if m[k] == nil {
			m[k] = v
		}
		return nil
	}
}

// Atoi - param is change to Int
func Atoi(k string) PCM {
	return func(m Model) error {
		v, ok := m[k].(string)
		if !ok {
			return NewErrParamWrong(k, "type error")
		}
		n, e := strconv.Atoi(v)
		if e != nil {
			return NewErrParamWrong(k, "type error")
		}
		m[k] = n
		return nil
	}
}
