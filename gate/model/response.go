package model

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/micro/go-log"
	microErr "github.com/micro/go-micro/errors"
)

// Response Model
type Response struct {
	c       echo.Context
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// NewResponse Builder
func NewResponse(c echo.Context, code int) *Response {
	if code == http.StatusOK {
		return &Response{
			c:       c,
			Code:    0,
			Message: "Success",
		}
	}
	return &Response{
		c:       c,
		Code:    -1,
		Message: http.StatusText(code),
	}
}

// Msg setter
func (r *Response) Msg(m string) *Response {
	r.Message = m
	return r
}

// MicroError Cleaner
func (r *Response) MicroError(e error) *Response {
	if e == nil {
		return r
	}
	if _, ok := e.(*microErr.Error); ok {
		log.Log(e)
		r.Message = ErrServiceCrash.Error()
		return r
	}
	return r.Msg(e.Error())
}

// Build response
func (r *Response) Build(d interface{}) error {
	r.Data = d
	return r.c.JSON(http.StatusOK, r)
}

// FreshJWT in header
func (r *Response) FreshJWT(jwt string) *Response {
	r.c.Response().Header().Set("Set-Token", jwt)
	return r
}
