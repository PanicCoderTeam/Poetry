package capi

import (
	"context"
	"poetry/src/pkg/log"
	"poetry/src/pkg/trpc/codec/capi_error"
	"time"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/plugin"
)

// var DefaultFramerBuilder = &CAPIFramer{}

func init() {
	p := &CAPI{}
	plugin.Register(pluginName, p)
	filter.Register(filterName, p.ServerFilter(), p.ClientFilter)
}

const (
	pluginName = "capi"
	pluginType = "api"
	filterName = "capi"
	capiFramer = "capi"
)

// CAPI 插件配置
type CAPI struct {
	AutoLogBody bool `yaml:"auto_log_body"`
}

// Validator 自动校验接口
type Validator interface {
	Validate() error
}

// ValidationError 自动校验接口
type ValidationError interface {
	Field() string
	Reason() string
}

// Setup 实现 plugin.Factory 接口。
func (c *CAPI) Setup(name string, dec plugin.Decoder) error {
	if err := dec.Decode(c); err != nil {
		return err
	}
	return nil
}

// Type 插件类型
func (c *CAPI) Type() string {
	return pluginType
}

// ServerFilter 服务端过滤器
func (c *CAPI) ServerFilter() filter.ServerFilter {
	return func(ctx context.Context, req any, next filter.ServerHandleFunc) (any, error) {
		requestId := ""
		head, ok := ctx.Value(http.ContextKeyHeader).(*http.Header)
		if ok {
			requestId = head.Request.Header.Get("RequestId")
		}
		begin := time.Now()
		log.WithContextFields(ctx, "RequestId", requestId)
		ip := trpc.GlobalConfig().Global.LocalIP
		log.WithContextFields(ctx, "server_ip", ip)
		rsp, e := next(ctx, req)
		cost := time.Since(begin) // 接受响应后计算耗时
		log.DebugContextEx(ctx, "req", req, "resp", rsp, "error", e, "cost", cost/1000)
		if e != nil {
			head.Response.Header().Add("ErrorCode", string(capi_error.ErrCode(e)))
		}
		return rsp, e
	}
}

func (c *CAPI) ClientFilter(ctx context.Context, req, rsp interface{}, next filter.ClientHandleFunc) error {
	requestId := ""
	head, ok := ctx.Value(http.ContextKeyHeader).(*http.Header)
	if ok {
		requestId = head.Request.Header.Get("RequestId")
	}
	begin := time.Now()
	log.WithContextFields(ctx, "RequestId", requestId)
	err := next(ctx, req, rsp)
	cost := time.Since(begin) // 接受响应后计算耗时
	if err != nil && ok {
		head.Response.Header().Add("ErrorCode", string(capi_error.ErrCode(err)))
	}
	log.DebugContextEx(ctx, "req", req, "resp", rsp, "error", err, "cost", cost/1000)
	return err
}

type CapiResponse struct {
	Response map[string]interface{}
}
