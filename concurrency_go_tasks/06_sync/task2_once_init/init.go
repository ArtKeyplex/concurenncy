package initonce

import "sync"

var (
	once        sync.Once
	initialized bool
)

// Init выполняет однократную инициализацию ресурса.
func Init() {
	// TODO: инициализировать ресурс через sync.Once
	once.Do(func() {
		initialized = true
	})
}

// Initialized возвращает, был ли инициализирован ресурс.
func Initialized() bool {
	// TODO: вернуть признак инициализации

	return initialized
}
