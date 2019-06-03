package router

import (
	"github.com/ddosakura/starmap/gate/router/auth"
	"github.com/labstack/echo"
)

// Bind router
func Bind(e *echo.Echo) {
	auth.Bind(e)
}
