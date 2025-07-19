package filter

import (
	"context"
	shttp "net/http"
	"poetry/src/pkg/basic"
	"poetry/src/pkg/log"
	"poetry/src/pkg/trpc/auth"
	"poetry/src/pkg/trpc/codec/capi_error"
	"strings"
	"time"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/http"
)

type AuthUserInfo struct {
	UserId   string
	UserName string
}

// Filter 实现filter接口
func Filter(ctx context.Context, req interface{}, next filter.ServerHandleFunc) (rsp interface{}, err error) {

	// 获取当前请求路径
	msg := trpc.Message(ctx)

	// 获取当前请求路径
	serviceName := msg.ServerRPCName()
	// 免拦截路径
	if strings.HasPrefix(serviceName, "/poetry/Login") || strings.HasPrefix(serviceName, "/poetry/CreateUser") {
		return next(ctx, req)
	}
	head, ok := ctx.Value(http.ContextKeyHeader).(*http.Header)
	log.DebugContextEx(ctx, "header", head.Request.Header, "headerType", head.Request.Method)
	if head.Request.Method == shttp.MethodOptions {
		return next(ctx, req)
	}
	authHeader := ""
	if ok {
		authHeader = head.Request.Header.Get("Authorization")
	}
	if len(authHeader) == 0 {
		return nil, capi_error.NewErr(capi_error.REQUEST_NOT_AUTH_CODE, "missing Authorization header")
	}
	claims, err := CheckToken(authHeader)
	if err != nil {
		return nil, err
	}
	// 将用户ID存入ctx
	newCtx := context.WithValue(ctx, basic.ClaimsKeyVal, AuthUserInfo{
		UserId:   claims.UserID,
		UserName: claims.UserName,
	})
	log.DebugContextEx(ctx, "auth push in ctx", newCtx.Value(basic.ClaimsKeyVal))
	return next(newCtx, req)
}

func CheckToken(authHeader string) (*auth.Claims, error) {
	// 校验格式
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, capi_error.NewError(capi_error.REQUEST_NOT_AUTH_CODE, "Token 格式无效", nil)
	}
	token := tokenParts[1]
	if string(token) == "" {
		return nil, capi_error.NewErr(capi_error.REQUEST_NOT_AUTH_CODE, "missing token")
	}

	// 验证token
	claims, err := auth.ParseToken(string(token))
	if err != nil {
		return nil, capi_error.NewErr(capi_error.REQUEST_NOT_AUTH_CODE, "invalid token")
	}
	if claims.ExpiresAt.Time.Unix() < time.Now().Unix() {
		return nil, capi_error.NewErr(capi_error.REQUEST_NOT_AUTH_CODE, "token expired")
	}
	return claims, nil
}
func ClientAuthFilter(ctx context.Context, req, rsp any, next filter.ClientHandleFunc) error {
	msg := trpc.Message(ctx)
	// 获取当前请求路径
	serviceName := msg.ClientRPCName()

	// 从上下文获取鉴权信息（如用户ID）
	authUserInfo, ok := ctx.Value(basic.ClaimsKeyVal).(AuthUserInfo)
	if !ok {
		// 免拦截路径
		if strings.HasPrefix(serviceName, "/poetry/Login") || strings.HasPrefix(serviceName, "/poetry/CreateUser") {
			return next(ctx, req, rsp)
		}
		return capi_error.NewErr(capi_error.REQUEST_NOT_AUTH_CODE, "claims not found in context")
	}

	md := msg.ClientMetaData()
	md["x-user-id"] = []byte(string(authUserInfo.UserId)) // 自定义透传字段
	head, ok := ctx.Value(http.ContextKeyHeader).(*http.Header)
	if ok {
		authHeader := head.Request.Header.Get("Authorization")
		if len(authHeader) > 0 {
			md["Authorization"] = []byte(authHeader)
		}
	}
	msg.WithClientMetaData(md)
	return next(ctx, req, rsp)
}

// Register 注册拦截器
func init() {
	filter.Register(basic.AuthFilterName, Filter, ClientAuthFilter)
}
