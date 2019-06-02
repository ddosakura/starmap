# RESTful

RESTful tools based on Middleware for micro api

> Warning:
> This is not a framework, because it is dependent on:
> + auth "github.com/ddosakura/starmap/srv/auth/proto" (auth service)
> + api "github.com/micro/go-api/proto" (not common api interface)
> + "github.com/micro/go-micro/errors"
> + client "github.com/micro/go-micro/client"

## 内置中间件使用顺序

### 认证顺序 (必须)

+ LoadAuthService
+ JWTCheck
+ RoleCheck/PermCheck/SuperRole

### 参数解析顺序 (建议)

+ ParamCheck
+ ParamAutoLoad (放在 ParamCheck 前面将无法使用默认值、重名名)

### 总体使用顺序 (建议)

+ LoadAuthService
+ JWTCheck
    + 须要请求认证服务，但在 RESTful 中的所有请求大都需要检查
+ ParamCheck
+ ParamAutoLoad (放在 ParamCheck 前面将无法使用默认值、重名名)
+ RoleCheck/PermCheck/SuperRole
    + 须要请求认证服务
    + RoleCheck 可能需要用 ParamAutoLoad 加载参数
