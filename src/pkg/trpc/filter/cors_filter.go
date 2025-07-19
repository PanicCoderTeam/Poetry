package filter

import (
	"context"
	"net/http"
	"poetry/src/pkg/log"

	"trpc.group/trpc-go/trpc-go/filter"
	thttp "trpc.group/trpc-go/trpc-go/http"
)

const CorsName = "cors_filter"

// Register 注册拦截器
func init() {
	filter.Register(CorsName, corsMiddleware, nil)
}

// 自定义 CORS 中间件
func corsMiddleware(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (interface{}, error) {
	// 获取 HTTP 响应对象
	head, ok := ctx.Value(thttp.ContextKeyHeader).(*thttp.Header)
	if ok {
		log.DebugContextEx(ctx, "cors middleware", "add cors header")
		// 设置 CORS 头
		head.Response.Header().Add("Access-Control-Allow-Origin", "*") // 允许前端源
		head.Response.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		head.Response.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, RequestId")
		head.Response.Header().Add("Access-Control-Allow-Credentials", "true") // 若需携带 Cookie
	}
	// 处理 OPTIONS 预检请求
	if head.Request.Method == http.MethodOptions {
		head.Response.WriteHeader(http.StatusNoContent) // 返回 204
		return nil, nil
	}

	// 继续处理业务逻辑
	return next(ctx, req)
}
