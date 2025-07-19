package log

import (
	"sync"
)

type LogManager struct {
	Logs    map[uint64]string
	LogLock sync.RWMutex
}

func NewLogManger() *LogManager {
	cm := &LogManager{
		Logs: make(map[uint64]string),
	}
	return cm
}
func (cm *LogManager) Add(id uint64, value string) {
	cm.LogLock.Lock()
	defer cm.LogLock.Unlock()
	cm.Logs[id] = value
}
func (cm *LogManager) Remove(id uint64) {
	cm.LogLock.Lock()
	defer cm.LogLock.Unlock()
	delete(cm.Logs, id)
}
func (cm *LogManager) Load(id uint64) (string, bool) {
	cm.LogLock.RLock()
	defer cm.LogLock.RUnlock()
	conn, ok := cm.Logs[id]
	if !ok {
		return "", false
	}
	return conn, true
}
func (cm *LogManager) Len() int {
	return len(cm.Logs)
}
func (cm *LogManager) Clean() {
	cm.LogLock.Lock()
	defer cm.LogLock.Unlock()
	for key := range cm.Logs {
		delete(cm.Logs, key)
	}
}
