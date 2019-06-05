package auth

import (
	"context"
	"net/http"

	mAuth "github.com/ddosakura/starmap/gate/middleware/auth"
	"github.com/ddosakura/starmap/gate/middleware/client"
	"github.com/ddosakura/starmap/gate/model"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/labstack/echo"
)

// API - 登录
// -> username, password
// <- userinfo
// !! set-token
func login(c echo.Context) error {
	ut, e := client.AuthUser(c).
		Login(context.Background(),
			c.Get(kuap).(*auth.UserAuth),
		)
	if e != nil {
		return model.NewResponse(c, http.StatusBadRequest).
			MicroError(e).
			Build(nil)
	}
	return model.NewResponse(c, http.StatusOK).
		FreshJWT(ut.Token).
		Build(ut.User)
}

// API - 注册
// -> username, password
// <- userinfo
// !! set-token
func register(c echo.Context) error {
	ut, e := client.AuthUser(c).
		Register(context.Background(),
			c.Get(kuap).(*auth.UserAuth),
		)
	if e != nil {
		return model.NewResponse(c, http.StatusBadRequest).
			MicroError(e).
			Build(nil)
	}
	return model.NewResponse(c, http.StatusOK).
		FreshJWT(ut.Token).
		Build(ut.User)
}

// API - 自身信息
// <- userinfo
// !! set-token
func info(c echo.Context) error {
	token := mAuth.GetJWT(c)
	ut, e := client.AuthUser(c).
		Info(context.Background(), &auth.UserToken{
			Token: token,
		})
	if e != nil {
		return model.NewResponse(c, http.StatusBadRequest).
			MicroError(e).
			Build(nil)
	}
	return model.NewResponse(c, http.StatusOK).
		FreshJWT(ut.Token).
		Build(ut.User)
}

// API - 更改自身信息
// -> pass/userinfo (optional, pass have priority)
// <- userinfo
// !! set-token
func update(c echo.Context) error {
	jwt := mAuth.JWT(c)
	u := c.Get(kuap).(*auth.UserToken)

	if u.Auth.Password == "" {
		u.User.UUID = jwt.User.UUID
		u.Auth = nil
	} else {
		u.Auth.ID = jwt.User.UUID
		u.User = nil
	}
	ut, e := client.AuthUser(c).Change(context.Background(), u)
	if e != nil {
		return model.NewResponse(c, http.StatusBadRequest).
			MicroError(e).
			Build(nil)
	}
	return model.NewResponse(c, http.StatusOK).
		FreshJWT(ut.Token).
		Build(ut.User)
}
