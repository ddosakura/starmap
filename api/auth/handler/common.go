package handler

import (
	"github.com/ddosakura/starmap/api/common"
	"github.com/ddosakura/starmap/api/rest"
)

var (
	autoLoadAuthService = rest.LoadAuthService(common.AuthUserFromContext)
)
