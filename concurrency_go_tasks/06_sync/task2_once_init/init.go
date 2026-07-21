package initonce

import "sync"

var (
	once        sync.Once
	initialized bool
)

// Init выполняет однократную инициализацию ресурса.
func Init() {
	once.Do(func() {
		initialized = true
	})
	// TODO: инициализировать ресурс через sync.Once
}

// Initialized возвращает, был ли инициализирован ресурс.
func Initialized() bool {
	// TODO: вернуть признак инициализации
	if initialized {
		return true
	} else {
		return false
	}
}
