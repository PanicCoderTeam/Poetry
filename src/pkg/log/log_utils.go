package log

import (
	"context"
	"encoding/json"

	"trpc.group/trpc-go/trpc-go/log"
)

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func WithContextFields(ctx context.Context, fields ...string) {
	log.WithContextFields(ctx, fields...)
}

// 调试
func DebugContextEx(ctx context.Context, strInterface ...interface{}) {
	strInterfaceParamMap := map[string]interface{}{}
	for i := 0; i < len(strInterface); i += 2 {
		if val, ok := strInterface[i].(string); !ok {
			strInterfaceParamMap["unKnowParam"] = strInterface[i]
			if i+1 >= len(strInterface) {
				continue
			}
			strInterfaceParamMap["unKnowParam"] = strInterface[i+1]
		} else {
			if i+1 >= len(strInterface) {
				continue
			}
			strInterfaceParamMap[val] = strInterface[i+1]
		}
	}
	logVal, _ := json.Marshal(strInterfaceParamMap)
	log.DebugContext(ctx, string(logVal))
}

// 错误信息
func ErrorContextEx(ctx context.Context, strInterface ...interface{}) {
	strInterfaceParamMap := map[string]interface{}{}
	for i := 0; i < len(strInterface); i += 2 {
		if val, ok := strInterface[i].(string); !ok {
			strInterfaceParamMap["unKnowParam"] = strInterface[i]
			if i+1 >= len(strInterface) {
				continue
			}
			strInterfaceParamMap["unKnowParam"] = strInterface[i+1]
		} else {
			if i+1 >= len(strInterface) {
				continue
			}
			strInterfaceParamMap[val] = strInterface[i+1]
		}
	}
	logVal, _ := json.Marshal(strInterfaceParamMap)
	log.ErrorContext(ctx, string(logVal))
}
func InfoContextEx(ctx context.Context, strInterface ...interface{}) {
	strInterfaceParamMap := map[string]interface{}{}
	for i := 0; i < len(strInterface); i += 2 {
		if val, ok := strInterface[i].(string); !ok {
			strInterfaceParamMap["unKnowParam"] = strInterface[i]
			if i+1 >= len(strInterface) {
				continue
			}
			strInterfaceParamMap["unKnowParam"] = strInterface[i+1]
		} else {
			if i+1 >= len(strInterface) {
				continue
			}
			strInterfaceParamMap[val] = strInterface[i+1]
		}
	}
	logVal, _ := json.Marshal(strInterfaceParamMap)
	log.InfoContext(ctx, string(logVal))
}
