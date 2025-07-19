package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"poetry/src/pkg/log"
	"time"
)

// DoRequest 云api调用封装
// response是传结构体的地址，传nil表示忽略回包。注意：只检查通信是否正常，不检查内容
func DoRequest(ctx context.Context, url string, header map[string]string, req interface{}, response interface{}) (
	int, error,
) {
	client := &http.Client{}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return -1, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return -2, err
	}
	// httpReq.
	httpReq = httpReq.WithContext(ctx)
	httpReq.Header.Set("Content-Type", "application/json")
	if len(header) > 0 {
		for k, v := range header {
			httpReq.Header.Set(k, v)
		}
	}
	start := time.Now()

	var httpRsp *http.Response
	httpRsp, err = client.Do(httpReq)
	if err != nil {
		return -3, err
	}
	defer httpRsp.Body.Close()
	var rspBody []byte
	rspBody, err = ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return -4, err
	}

	log.DebugEx(
		ctx, "DoRequest", "cost", time.Since(start), "url", url, "header", httpReq.Header, log.Json("req", reqBody),
		log.Json("rsp", rspBody),
	)
	if len(rspBody) != 0 {
		err = json.Unmarshal(rspBody, response)
		if err != nil {
			return -9, errors.New("http rsp json format err:" + err.Error())
		}
	}
	return 0, nil
}
