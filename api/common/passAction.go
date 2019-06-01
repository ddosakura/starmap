package common

type passAction struct {
	s *RESTful
}

func (a *passAction) Role(rules []string, logical Logical) Action {
	return a
}

func (a *passAction) Permission(rules []string, logical Logical) Action {
	return a
}

func (a *passAction) Check(k string, multi bool, defautlV []string) Action {
	return a
}

func (a *passAction) Do(fn func(*RESTful) (interface{}, error)) *RESTful {
	return a.s
}
