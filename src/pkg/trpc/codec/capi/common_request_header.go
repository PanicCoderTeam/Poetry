package capi

import (
	"context"
	"poetry/src/pkg/trpc/codec/capi_error"

	"trpc.group/trpc-go/trpc-go/http"
)

const (
	HEADER_REQUEST_ID = "RequestId"
)

type Header struct {
	RequestId string
	Language  string
	ErrorCode capi_error.ErrorCode
}

func GetRequestId(ctx context.Context) string {
	head, ok := ctx.Value(http.ContextKeyHeader).(*http.Header)
	if ok {
		requestId := head.Request.Header.Get("RequestId")
		return requestId
	}
	return ""
}
