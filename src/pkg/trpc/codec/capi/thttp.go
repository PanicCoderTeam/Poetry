package capi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"poetry/src/pkg/trpc/codec/capi_error"

	"trpc.group/trpc-go/trpc-go/errs"
	thttp "trpc.group/trpc-go/trpc-go/http"
)

func init() {
	thttp.DefaultServerCodec.ErrHandler = func(w http.ResponseWriter, r *http.Request, e *errs.Error) {
		var buf bytes.Buffer
		buf.Reset()
		requestId := r.Header.Get("RequestId")
		language := r.Header.Get("Language")
		capiError := capi_error.INTERNAL_ERROR_CODE.ErrInfo(language)
		if e != nil {
			// log.DebugContextEx(context.Background(), "debug error", w.Header(), "requestId", requestId)
			errCode := capi_error.ErrCode(e)
			if w.Header() != nil && len(w.Header().Get("Errorcode")) > 0 {
				errCode = capi_error.ErrorCode(w.Header().Get("Errorcode"))
			}
			// log.DebugContextEx(context.Background(), "debug error", e, "requestId", requestId)
			// log.DebugContextEx(context.Background(), "errCode", errCode, "requestId", requestId)
			capiError = errCode.ErrInfo(language)
		}
		buf.WriteString("{\"Response\":{\"RequestId\":\"" + requestId + "\",\"Error\": ")
		errMsg, err := json.Marshal(capiError)
		if err == nil {
			buf.WriteString(string(errMsg))
		}
		buf.WriteString("}}")
		w.Write(buf.Bytes())
	}
	thttp.DefaultServerCodec.RspHandler = func(w http.ResponseWriter, r *http.Request, rspBody []byte) error {
		if len(rspBody) == 0 {
			return nil
		}
		var buf bytes.Buffer
		buf.Reset()
		requestId := r.Header.Get("RequestId")
		buf.WriteString("{\"Response\":{\"RequestId\":\"" + requestId + "\",\"Data\":" + string(rspBody) + "}}")
		if _, err := w.Write(buf.Bytes()); err != nil {
			return fmt.Errorf("http write response error: %s", err.Error())
		}
		return nil
	}
}
