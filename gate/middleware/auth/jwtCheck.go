package auth

import (
	"context"
	"net/http"

	"github.com/ddosakura/starmap/gate/middleware/client"
	"github.com/ddosakura/starmap/gate/model"
	auth "github.com/ddosakura/starmap/srv/auth/proto"
	"github.com/labstack/echo"
)

const (
	jwtToken = "jwt"
)

// GetJWT from header
func GetJWT(c echo.Context) string {
	ps := c.Request().Header["Token"]
	if ps == nil || len(ps) == 0 {
		return ""
	}
	return ps[0]
}

// JWTCheck for Client
func JWTCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, e := client.AuthUser(c).Check(context.Background(), &auth.UserToken{
			Token: GetJWT(c),
		})

		if e != nil {
			return model.NewResponse(c, http.StatusBadRequest).
				MicroError(e).
				Build(nil)
		}
		c.Set(jwtToken, token)
		if token.Token != "" {
			model.NewResponse(c, http.StatusOK).FreshJWT(token.Token)
		}
		return next(c)
	}
}

// JWT getter
func JWT(c echo.Context) *auth.UserToken {
	return c.Get(jwtToken).(*auth.UserToken)
}
