package echo

// 为了尽量做一个框架无关的参数验证库，echo 相关的代码会放在这里

import (
	"net/http"
	"net/url"

	"github.com/ddosakura/starmap/gate/check"
	"github.com/labstack/echo"
)

func init() {
	check.Build = ParamCheckBuilder
}

// ParamCheckBuilder Export
func ParamCheckBuilder(i interface{}) check.RuleBuilder {
	c, ok := i.(echo.Context)
	b := &check.BaseRuleBuilder{}
	if !ok {
		b.Err = check.ErrBuildConfig
		return b
	}
	switch c.Request().Method {
	case http.MethodGet, http.MethodDelete:
		b.Params = c.QueryParams()
		//b.Param = c.QueryParam
	case http.MethodPost, http.MethodPut:
		params, e := c.FormParams()
		if e != nil {
			b.Err = e
		} else {
			b.Params = params
			//b.Param = c.FormValue
		}
	default:
		b.Err = check.ErrNotRESTfulRequest
	}
	return b
}

// PathParamCheckBuilder Export
func PathParamCheckBuilder(i interface{}) check.RuleBuilder {
	c, ok := i.(echo.Context)
	b := &check.BaseRuleBuilder{}
	if !ok {
		b.Err = check.ErrBuildConfig
		return b
	}
	b.Params = make(url.Values)
	for _, k := range c.ParamNames() {
		b.Params.Set(k, c.Param(k))
	}
	//b.Param = c.Param
	return b
}
