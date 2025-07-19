package basic

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

/*  Build-In Keys */

type contextKey struct{}

/* Constants */

var ctxData contextKey
var errNotInitialized = errors.New("ctx not initialized")
var errDataNotStored = errors.New("data not stored")

/* Context Data Manipulation */

func InitCtx(rawCtx context.Context) context.Context {
	if rawCtx.Value(ctxData) != nil {
		return rawCtx
	} else {
		return context.WithValue(rawCtx, ctxData, &sync.Map{})
	}
}

func Background() context.Context {
	return context.WithValue(context.Background(), ctxData, &sync.Map{})
}

func BackgroundWithData(rawCtx context.Context) context.Context {
	v := AllValues(rawCtx)
	if v != nil {
		return context.WithValue(context.Background(), ctxData, v)
	} else {
		return context.Background() // 避免interface不为nil
	}
}

func AllValues(ctx context.Context) *sync.Map {
	if data, ok := ctx.Value(ctxData).(*sync.Map); !ok {
		return nil
	} else {
		return data
	}
}

func StoreVal(ctx context.Context, key interface{}, val interface{}) error {
	if data, ok := ctx.Value(ctxData).(*sync.Map); !ok {
		return errNotInitialized
	} else {
		data.Store(key, val)
		return nil
	}
}

func StoreUniqueVal(ctx context.Context, key interface{}, val interface{}) error {
	if data, ok := ctx.Value(ctxData).(*sync.Map); !ok {
		return errNotInitialized
	} else {
		_, loaded := data.LoadOrStore(key, val)
		if loaded {
			return fmt.Errorf("key:%+v already exists", key)
		}
		return nil
	}
}

func LoadVal(ctx context.Context, key interface{}) (interface{}, bool) {
	if data, ok := ctx.Value(ctxData).(*sync.Map); !ok {
		return nil, false
	} else {
		return data.Load(key)
	}
}

func GetVal(ctx context.Context, key interface{}) interface{} {
	if data, ok := ctx.Value(ctxData).(*sync.Map); !ok {
		return nil
	} else {
		v, _ := data.Load(key)
		return v
	}
}

func MustStoreVal(ctx context.Context, key interface{}, val interface{}) {
	if data, ok := ctx.Value(ctxData).(*sync.Map); !ok {
		panic(errNotInitialized)
	} else {
		data.Store(key, val)
	}
}

func MustLoadVal(ctx context.Context, key interface{}) interface{} {
	if data, ok := ctx.Value(ctxData).(*sync.Map); !ok {
		panic(errNotInitialized)
	} else if v, ok := data.Load(key); !ok {
		panic(errDataNotStored)
	} else {
		return v
	}
}
