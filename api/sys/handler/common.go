package handler

import (
	"context"

	"github.com/ddosakura/starmap/api/common"
	"github.com/ddosakura/starmap/api/rest"
)

var (
	autoLoadAuthService = rest.LoadAuthService(func(ctx context.Context) (rest.AuthService, bool) {
		return common.AuthUserFromContext(ctx)
	})
)
