package fileflow

import "sync"

// 全局注册表：用于业务侧在初始化阶段注册 Processor/Middleware。
var registry struct {
	mu            sync.RWMutex
	processors    []Processor
	middlewares   []Middleware
	postProcessor PostProcessor
}

// RegisterProcessor 注册 Processor。
func RegisterProcessor(proc Processor) {
	if proc == nil {
		return
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.processors = append(registry.processors, proc)
}

// RegisterMiddleware 注册 Middleware。
func RegisterMiddleware(m Middleware) {
	if m == nil {
		return
	}
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.middlewares = append(registry.middlewares, m)
}

// RegisterPostProcessor 注册全局 PostProcessor（后注册会覆盖前一个）。
func RegisterPostProcessor(pp PostProcessor) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.postProcessor = pp
}

// SnapshotRegistered 返回当前注册项快照，避免调用方直接持有内部切片。
func SnapshotRegistered() (processors []Processor, middlewares []Middleware, postProcessor PostProcessor) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	processors = append(processors, registry.processors...)
	middlewares = append(middlewares, registry.middlewares...)
	postProcessor = registry.postProcessor
	return
}

// HasRegisteredProcessor 返回是否至少注册了一个 Processor。
func HasRegisteredProcessor() bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	return len(registry.processors) > 0
}
