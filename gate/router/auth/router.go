package auth

import (
	"net/http"

	"github.com/ddosakura/starmap/gate/check"
	mAuth "github.com/ddosakura/starmap/gate/middleware/auth"
	"github.com/ddosakura/starmap/gate/model"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/labstack/echo"
)

// Bind router
func Bind(e *echo.Echo) {
	g := e.Group("/auth")

	g.POST("/login", login, userAuthParam).Name = "auth.login"
	g.POST("/register", register, userAuthParam).Name = "auth.register"

	// didn't JWTCheck, beack Srv.Auth.Info will also check it
	g.POST("/info", info).Name = "auth.info"

	g.POST("/update", update, userAuthPassWithUserInfo, mAuth.JWTCheck).Name = "auth.update"
}

const (
	kuap = "uap" // key of userAuthParam
)

func userAuthParam(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		a := new(auth.UserAuth)
		if e := check.
			Build(c).
			Rules(check.M("user", "pass"),
				check.Rename("id", ""), // 过滤 id 参数
				check.Rename("user", "username"),
				check.Rename("pass", "password"),
			).Load(a); e != nil {
			return model.NewResponse(c, http.StatusBadRequest).
				Msg(e.Error()).
				Build(nil)
		}
		c.Set(kuap, a)
		return next(c)
	}
}

func userAuthPassWithUserInfo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(auth.UserToken)
		u.Auth = new(auth.UserAuth)
		u.User = new(auth.UserInfo)
		if e := check.
			Build(c).
			Rules(nil,
				check.Rename("id", ""),
				check.Rename("user", ""),
				check.Rename("pass", "password"),
			).Load(u.Auth, u.User); e != nil {
			return model.NewResponse(c, http.StatusBadRequest).
				Msg(e.Error()).
				Build(nil)
		}
		c.Set(kuap, u)
		return next(c)
	}
}
