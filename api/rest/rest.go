package rest

import (
	"context"

	api "github.com/micro/go-api/proto"
)

// RESTfulType for API
type RESTfulType int

// RESTfulType(s)
const (
	POST   RESTfulType = 1 << iota // INSERT
	DELETE                         // DELETE
	GET                            // INSERT
	PUT                            // UPDATE
)

// Configurable raw
var (
	SrvName = "starmap.api"
)

// dict
var (
	RESTfulTypeDict = map[string]RESTfulType{
		"GET":    GET,
		"get":    GET,
		"POST":   POST,
		"post":   POST,
		"delete": DELETE,
		"DELETE": DELETE,
		"put":    PUT,
		"PUT":    PUT,
	}
)

// RESTful API
type RESTful interface {
	// flow
	Action(RESTfulType) RESTful // build child
	Chain(Middleware) RESTful   // use Middleware
	Done() RESTful              // backto father
	Final() error               // error response (nil means success)
}

// REST Builder
func REST(ctx context.Context, req *api.Request, res *api.Response) RESTful {
	s := &Flow{
		Ctx:  ctx,
		Req:  req,
		Res:  res,
		Rest: RESTfulTypeDict[req.Method],

		ctxForMiddle: context.Background(),
		final:        false,
	}
	switch s.Rest {
	case GET, DELETE:
		s.Params = req.Get
	case POST, PUT:
		s.Params = req.Post
	default:
		s.Params = make(map[string]*api.Pair)
	}
	return s
}
