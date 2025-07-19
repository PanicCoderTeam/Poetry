package log

import (
	"encoding/json"

	"go.uber.org/zap"
)

func Json(key string, data []byte) zap.Field {
	if json.Valid(data) {
		return zap.Reflect(key, json.RawMessage(data))
	} else {
		return zap.ByteString(key, data)
	}
}
