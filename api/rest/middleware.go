package rest

import (
	"context"
	"encoding/json"
)

// Middleware for RESTful
type Middleware func(context.Context, *Flow) error // if it returns a error, flow will finish

// --- APIs ---

// Success API
func (s *Flow) Success(d interface{}) error {
	if s.final {
		return s.err
	}
	s.final = true

	b, e := json.Marshal(&struct {
		Code   int         `json:"code"`
		Detail string      `json:"detail"`
		Data   interface{} `json:"data"`
	}{
		Code:   0,
		Detail: "Success",
		Data:   d,
	})
	if e != nil {
		s.err = e
	} else {
		s.Res.StatusCode = 200
		s.Res.Body = string(b)
	}
	return s.err
}

// FreshJWT API
func (s *Flow) FreshJWT(jwt string) {
	SetJWT(s.Res, jwt)
}
