package basic

import (
	"context"
	"fmt"
	"sync"
)

// ErrorCollector 用于在多个 goroutine 之间收集错误信息
type ErrorCollector interface {
	Put(...interface{})
	PutWithCancel(values ...interface{})

	HasError() bool
	GetAll() []error
}

type errorCollectorImpl struct {
	sync.RWMutex

	buffer       []error
	cancelHelper context.CancelFunc // 用来取消父context
}

// NewErrorCollector ...
// 不需要使用 PutWithCancel 接口时，cancelHelper 可传 nil
func NewErrorCollector(cancelHelper context.CancelFunc) ErrorCollector {
	return &errorCollectorImpl{cancelHelper: cancelHelper}
}

// Put ...
func (e *errorCollectorImpl) Put(values ...interface{}) {
	e.Lock()
	defer e.Unlock()

	e.buffer = append(e.buffer, wrapToError(values)...)
}

// PutWithCancel ...
func (e *errorCollectorImpl) PutWithCancel(values ...interface{}) {
	e.Put(values...)
	if e.cancelHelper != nil {
		e.cancelHelper()
	}
}

// HasError ...
func (e *errorCollectorImpl) HasError() bool {
	e.RLock()
	defer e.RUnlock()

	return len(e.buffer) > 0
}

// GetAll ...
func (e *errorCollectorImpl) GetAll() []error {
	e.RLock()
	defer e.RUnlock()

	return e.buffer
}

/* Helper */

func wrapToError(values []interface{}) []error {
	var result []error
	for _, v := range values {
		switch x := v.(type) {
		case error:
			result = append(result, x)
		case []error:
			result = append(result, x...)
		default:
			result = append(result, fmt.Errorf("%+v", x))
		}
	}
	return result
}
