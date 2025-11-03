package limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	t.Parallel()
	l := NewLimiter()
	defer l.Stop()

	var allowed int32
	for i := 0; i < 10; i++ {
		if l.Allow() {
			atomic.AddInt32(&allowed, 1)
		}
	}
	if allowed > 5 {
		t.Fatalf("нарушение лимита: разрешено более 5 событий сразу (%d). Лимитер должен позволять максимум 5 токенов в начальный момент, что указывает на проблему с инициализацией или механизмом ограничения", allowed)
	}
	time.Sleep(1100 * time.Millisecond)
	for i := 0; i < 5; i++ {
		if !l.Allow() {
			t.Fatalf("ожидалось, что токен %d будет разрешён после пополнения через 1100мс, но доступ был отклонён. Это указывает на проблему с механизмом пополнения токенов или таймером", i)
		}
	}
}

func TestNewLimiterNotNil(t *testing.T) {
	t.Parallel()
	l := NewLimiter()
	defer l.Stop()
	if l == nil {
		t.Fatal("ожидался непустой лимитер")
	}
}

func TestLimiterInitialTokens(t *testing.T) {
	t.Parallel()
	l := NewLimiter()
	defer l.Stop()

	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}

	if allowed != 5 {
		t.Fatalf("неверное количество начальных токенов: ожидали 5 токенов при инициализации лимитера, получили %d. Это может указывать на проблему с заполнением начального буфера токенов", allowed)
	}
}

func TestLimiterRefill(t *testing.T) {
	t.Parallel()
	l := NewLimiter()
	defer l.Stop()

	for i := 0; i < 5; i++ {
		if !l.Allow() {
			t.Fatalf("ожидалось, что токен %d из начальных 5 будет разрешён, но доступ был отклонён. Это указывает на проблему с начальными токенами или механизмом проверки доступности", i)
		}
	}

	if l.Allow() {
		t.Fatal("нарушение лимита: лимитер разрешил более 5 токенов сразу после потребления всех начальных токенов. Лимитер должен блокировать дополнительные запросы до пополнения токенов")
	}

	time.Sleep(250 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("ожидалось, что токен будет доступен после пополнения через 250мс (1 токен пополняется каждые 200мс), но доступ был отклонён. Это указывает на проблему с механизмом пополнения или таймером")
	}
}

func TestLimiterRateLimit(t *testing.T) {
	t.Parallel()
	l := NewLimiter()
	defer l.Stop()

	var allowed int32
	start := time.Now()

	for i := 0; i < 10; i++ {
		if l.Allow() {
			atomic.AddInt32(&allowed, 1)
		}
		time.Sleep(10 * time.Millisecond)
	}

	duration := time.Since(start)

	if atomic.LoadInt32(&allowed) > 6 {
		t.Fatalf("нарушение ограничения скорости: разрешено слишком много токенов (%d) за период %v. При скорости пополнения 1 токен в 200мс, за ~100мс должно быть разрешено максимум 5-6 токенов (5 начальных + 1 пополненный). Это указывает на проблему с механизмом ограничения скорости", allowed, duration)
	}
}

func TestLimiterSustainedRate(t *testing.T) {
	t.Parallel()
	l := NewLimiter()
	defer l.Stop()

	var allowed int32
	duration := 2 * time.Second
	start := time.Now()

	for time.Since(start) < duration {
		if l.Allow() {
			atomic.AddInt32(&allowed, 1)
		}
		time.Sleep(10 * time.Millisecond)
	}

	totalAllowed := atomic.LoadInt32(&allowed)
	expectedMin := int32(13)
	expectedMax := int32(15)

	if totalAllowed < expectedMin || totalAllowed > expectedMax {
		t.Errorf("неверная скорость пополнения токенов: разрешено %d токенов за %v, ожидалось между %d и %d. За 2 секунды должно быть доступно примерно 15 токенов (5 начальных + 10 пополненных через каждые 200мс). Это указывает на проблему с механизмом пополнения или таймером",
			totalAllowed, duration, expectedMin, expectedMax)
	}
}

func TestLimiterConcurrentAccess(t *testing.T) {
	t.Parallel()
	l := NewLimiter()
	defer l.Stop()

	const numGoroutines = 20
	const attemptsPerGoroutine = 10
	var allowed int32

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < attemptsPerGoroutine; j++ {
				if l.Allow() {
					atomic.AddInt32(&allowed, 1)
				}
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	totalAllowed := atomic.LoadInt32(&allowed)
	if totalAllowed > 20 {
		t.Errorf("нарушение лимита при параллельном доступе: разрешено %d токенов при одновременном обращении 20 горутин. Даже при параллельном доступе лимитер должен соблюдать ограничение в 5 токенов в секунду. Это указывает на проблему с потокобезопасностью механизма ограничения", totalAllowed)
	}
}

func TestLimiterMultipleInstances(t *testing.T) {
	t.Parallel()
	l1 := NewLimiter()
	defer l1.Stop()
	l2 := NewLimiter()
	defer l2.Stop()

	allowed1 := 0
	allowed2 := 0

	for i := 0; i < 10; i++ {
		if l1.Allow() {
			allowed1++
		}
		if l2.Allow() {
			allowed2++
		}
	}

	if allowed1 != 5 {
		t.Errorf("неверное количество токенов для лимитера 1: ожидали 5 начальных токенов, получили %d. Каждый экземпляр лимитера должен иметь свой независимый буфер токенов. Это указывает на проблему с изоляцией экземпляров", allowed1)
	}
	if allowed2 != 5 {
		t.Errorf("неверное количество токенов для лимитера 2: ожидали 5 начальных токенов, получили %d. Каждый экземпляр лимитера должен иметь свой независимый буфер токенов. Это указывает на проблему с изоляцией экземпляров", allowed2)
	}
}

func TestLimiterStop(t *testing.T) {
	t.Parallel()
	l := NewLimiter()

	for i := 0; i < 5; i++ {
		l.Allow()
	}

	l.Stop()

	time.Sleep(300 * time.Millisecond)
	if l.Allow() {
		t.Fatal("нарушение остановки лимитера: после вызова Stop() лимитер не должен разрешать новые токены даже после ожидания. Токены не должны пополняться после остановки. Это указывает на проблему с механизмом остановки или утечкой горутины пополнения")
	}
}

func TestLimiterStopMultipleTimes(t *testing.T) {
	t.Parallel()
	l := NewLimiter()
	l.Stop()
	l.Stop()
	l.Stop()
}

func BenchmarkLimiterAllow(b *testing.B) {
	l := NewLimiter()
	defer l.Stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Allow()
	}
}

func BenchmarkLimiterConcurrent(b *testing.B) {
	l := NewLimiter()
	defer l.Stop()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = l.Allow()
		}
	})
}
