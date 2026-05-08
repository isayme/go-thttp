# thttp

统一 HTTP 框架，适配多种路由库（gin/echo/chi/httprouter/gorilla-mux/net/http）。

## 项目结构

```
thttp.go          # App 主类，HTTP 服务入口
context.go        # Context 接口，实现请求上下文
router.go         # Router 接口定义及类型映射
route.go          # 路由解析（Static/Param/CatchAll）
group.go          # 路由分组
middleware.go     # 中间件机制
option.go         # 配置选项

router_*.go       # 各路由适配实现
- router_std.go        # net/http
- router_gin.go        # gin-gonic/gin
- router_echo.go       # labstack/echo
- router_chi.go        # go-chi/chi
- router_httprouter.go # julienschmidt/httprouter
- router_gorilla_mux.go

middleware/      # 内置中间件
- recovery.go
- logger.go
- requestid.go
- basic_auth.go
```

## 核心接口

```go
// Router 路由适配接口
type Router interface {
    Use(middleware ...MiddlewareFunc)
    Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc)
    Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool)
    FormatSegment(seg Segment) string
}

// Context 请求上下文
type Context interface {
    Request() *http.Request
    Response() http.ResponseWriter
    PathParam(name string) string
    QueryParam(name string) string
    JSON(code int, i interface{}) error
    // ...
}
```

## 路由切换

通过环境变量 `THTTP_ROUTER_TYPE` 或 `WithRouterType` 选项切换：

```go
// 环境变量
THTTP_ROUTER_TYPE=gin thttp.New()

// 代码
thttp.New(thttp.WithRouterType(thttp.RouterTypeGin))
```

## 路由模式

支持三种模式，通过 `FormatSegment` 适配不同框架的语法：

| 模式 | thttp 写法 | gin | echo | chi | std |
|------|-----------|-----|------|-----|-----|
| 静态 | `/users` | `/users` | `/users` | `/users` | `/users` |
| 参数 | `/users/:id` | `/users/:id` | `/users/:id` | `/users/{id}` | `/users/{id}` |
| 捕获 | `/users/*path` | `/users/*path` | `/users/*path` | `/users/*` | `/users/{path...}` |

## 注意事项

1. **多路由匹配差异**：不同框架对同一请求的路由匹配行为不同
   - gin: 支持多个路由匹配同一路径，按注册顺序执行
   - echo: 同样支持多路由
   - chi/std: 单一匹配

2. **参数获取方式**：各适配器通过 `PathParamsFunc` 转换框架原生的参数获取方式

3. **fakeResponseWriter**：gin 适配使用 fakeResponseWriter 避免框架自动写响应