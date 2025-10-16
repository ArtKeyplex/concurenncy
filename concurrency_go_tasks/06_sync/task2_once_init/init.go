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
}

// Initialized возвращает, был ли инициализирован ресурс.
func Initialized() bool {
	return initialized
}
