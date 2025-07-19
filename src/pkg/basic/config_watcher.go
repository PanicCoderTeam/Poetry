package basic

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	defaultFrequency = 5 * time.Second
)

/* API */

type (
	// ConfigWatcher 用于读取配置文件并监视文件内容的变化
	ConfigWatcher interface {
		// 添加监视文件。tag 为配置的标志，用于取回数据。path 为配置文件路径，data 是配置所对应的数据结构的指针（默认文件内容是 JSON）
		// options 详见各个配置的注释
		Watch(tag, path string, data interface{}, options ...WatchOption) error
		// 根据 tag 取回配置数据
		Retrieve(tag string) interface{}
		// 停止监视文件变化，释放资源
		Release()
	}

	ConfigParseFunc     func([]byte, interface{}) error
	ConfigWatcherOption func(*configWatcher)
	WatchOption         func(*watchedConfig)
)

/* Option Impl */

// LoggerOpt 修改 Watcher 的 Logger，不提供时更新文件失败等异常日志不会打印
func LoggerOpt(logger *zap.Logger) ConfigWatcherOption {
	return func(w *configWatcher) {
		w.SugaredLogger = logger.Sugar()
	}
}

// FrequencyOpt 修改 Watcher 的监控频率，单位为秒，默认为 5 秒
func FrequencyOpt(seconds uint) ConfigWatcherOption {
	return func(w *configWatcher) {
		w.frequency = time.Second * time.Duration(seconds)
	}
}

// ParseOpt 指定解析数据时使用的方法
func ParseOpt(f ConfigParseFunc) WatchOption {
	return func(c *watchedConfig) {
		c.parseFunc = f
	}
}

/* Watcher Impl */

type watchedConfig struct {
	sync.RWMutex

	path          string
	data          interface{}
	parseFunc     ConfigParseFunc
	modTimeRecord time.Time
}

func (w *watchedConfig) SetData(bytes []byte) error {
	v := reflect.ValueOf(w.data)
	nv := reflect.New(v.Type().Elem())
	err := w.parseFunc(bytes, nv.Interface())
	if err != nil {
		return err
	}

	w.Lock()
	w.data = nv.Interface()
	w.Unlock()
	return nil
}

func (w *watchedConfig) Data() interface{} {
	w.RLock()
	defer w.RUnlock()

	return w.data
}

func (w *watchedConfig) update() error {
	if info, err := os.Stat(w.path); err != nil {
		return fmt.Errorf("get file info failed: %v", err)

	} else if !w.modTimeRecord.IsZero() && w.modTimeRecord.Equal(info.ModTime()) {
		return nil // 无需更新

	} else if bytes, err := ioutil.ReadFile(w.path); err != nil {
		return fmt.Errorf("read file filed: %v", err)

	} else if err := w.SetData(bytes); err != nil {
		return fmt.Errorf("parse config failed: %v", err)

	} else {
		w.modTimeRecord = info.ModTime()
		return nil
	}
}

type configWatcher struct {
	sync.RWMutex
	*zap.SugaredLogger

	table     map[string]*watchedConfig // tag to config
	frequency time.Duration
	cancel    context.CancelFunc
}

func NewConfigWatcher(option ...ConfigWatcherOption) ConfigWatcher {
	ctx, cancel := context.WithCancel(context.Background())
	w := &configWatcher{table: make(map[string]*watchedConfig), frequency: defaultFrequency, cancel: cancel}
	for _, optFunc := range option {
		optFunc(w)
	}
	if w.SugaredLogger == nil {
		w.SugaredLogger = zap.NewNop().Sugar()
	}

	go w.updateLoop(ctx)
	return w
}

func (c *configWatcher) Watch(tag, path string, data interface{}, options ...WatchOption) error {
	c.RLock()
	_, ok := c.table[tag]
	c.RUnlock()
	if ok {
		return fmt.Errorf("duplicated tag: %s", tag)
	}
	if v := reflect.ValueOf(data); v.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid config data: %T(%v)", data, data)
	}

	wc := &watchedConfig{path: path, data: data}
	for _, optFunc := range options {
		optFunc(wc)
	}
	if wc.parseFunc == nil {
		wc.parseFunc = json.Unmarshal
	}
	if err := wc.update(); err != nil {
		return err
	}

	c.Lock()
	if _, ok := c.table[tag]; ok {
		return fmt.Errorf("duplicated tag: %s", tag)
	}
	c.table[tag] = wc
	c.Unlock()
	return nil
}

func (c *configWatcher) Retrieve(tag string) interface{} {
	c.RLock()
	defer c.RUnlock()

	if wc, ok := c.table[tag]; !ok {
		c.Errorw("can't find config", "tag", tag)
		return nil
	} else {
		return wc.Data()
	}
}

func (c *configWatcher) Release() {
	c.cancel()
}

func (c *configWatcher) updateLoop(ctx context.Context) {
	c.Infow("config watcher started")
	ticker := time.NewTicker(c.frequency)
	for {
		select {
		case <-ticker.C:
			c.updateAll()
		case <-ctx.Done():
			ticker.Stop()
			c.Infow("config watcher stopped")
			return
		}
	}
}

func (c *configWatcher) updateAll() {
	c.RLock()
	defer c.RUnlock() // 这里锁住的临界资源是 c.table，所以用读锁
	defer func() {
		if err := recover(); err != nil {
			c.Errorw("panic occurred", "err", err)
		}
	}()

	for tag, wc := range c.table {
		if err := wc.update(); err != nil {
			c.Errorw("update config failed", "tag", tag, "path", wc.path, "err", err)
		}
	}
}
