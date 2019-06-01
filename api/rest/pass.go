package rest

type pass struct {
	f RESTful
}

func (s *pass) Action(RESTfulType) RESTful {
	return &pass{s}
}

func (s *pass) Chain(Middleware) RESTful {
	return s
}

func (s *pass) Done() RESTful {
	return s.f
}

func (s *pass) Final() error {
	return nil
}
