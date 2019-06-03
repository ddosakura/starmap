package auth

import (
	"context"
	"net/http"

	"github.com/ddosakura/starmap/gate/check"
	"github.com/ddosakura/starmap/gate/middleware/client"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/labstack/echo"
)

// Bind router
func Bind(e *echo.Echo) {
	g := e.Group("/auth")
	//g.Use(middleware.BasicAuth(func(username, password string) bool {
	//	if username == "joe" && password == "secret" {
	//		return true
	//	}
	//	return false
	//}))

	g.GET("/login", login).Name = "auth.login"
}

func login(c echo.Context) error {
	a := new(auth.UserAuth)
	if e := check.
		Build(c).
		Rules(check.M("user", "pass"),
			check.Rename("user", "username"),
			check.Rename("pass", "password"),
		).Load(a); e != nil {
		return c.JSON(http.StatusOK, struct {
			Code int
			Msg  string
		}{Code: -1, Msg: e.Error()})
	}

	ut, e := client.AuthUser(c).Login(context.Background(), a)
	if e != nil {
		return c.JSON(http.StatusOK, struct {
			Code int
			Msg  string
		}{Code: -1, Msg: e.Error()})
	}
	return c.JSON(http.StatusOK, struct {
		Code int
		Data interface{}
	}{Code: 0, Data: ut.User})
}
