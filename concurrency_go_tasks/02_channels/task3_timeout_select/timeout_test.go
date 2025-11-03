package timeout

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWorkTimeout(t *testing.T) {
	t.Parallel()
	err := Work(context.Background())
	if err == nil {
		t.Fatal("ожидалась ошибка таймаута: функция Work() должна возвращать ошибку ErrTimeout при работе с контекстом без дедлайна. При отсутствии таймаута функция должна завершиться с ошибкой. Это указывает на проблему с механизмом таймаута")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("неверный тип ошибки: ожидали ErrTimeout, получили %v (тип: %T). Функция должна возвращать специфическую ошибку таймаута, а не другую ошибку", err, err)
	}
}

func TestWorkCanceled(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Work(ctx)
	if err == nil {
		t.Fatal("ожидалась ошибка отмены контекста")
	}
	if !errors.Is(err, ErrCanceled) {
		t.Errorf("неверный тип ошибки: ожидали ErrCanceled для отменённого контекста, получили %v (тип: %T). При отмене контекста функция должна возвращать ошибку отмены, а не другую ошибку", err, err)
	}
}

func TestWorkTimeoutBeforeCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := Work(ctx)
	if err == nil {
		t.Fatal("ожидалась ошибка таймаута: функция Work() должна возвращать ошибку ErrTimeout даже при наличии контекста с возможностью отмены, если таймаут происходит раньше отмены. Это указывает на проблему с приоритетом обработки таймаута")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("неверный тип ошибки: ожидали ErrTimeout, получили %v (тип: %T). При таймауте функция должна возвращать ошибку таймаута, а не другую ошибку", err, err)
	}
}

func TestWorkCanceledContext(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := Work(ctx)
	if err == nil {
		t.Fatal("ожидалась ошибка от отменённого контекста: функция Work() должна возвращать ошибку при отмене контекста во время выполнения. Контекст отменяется через 50мс, что должно привести к ошибке ErrCanceled или ErrTimeout. Это указывает на проблему с обработкой отмены контекста")
	}
	if !errors.Is(err, ErrTimeout) && !errors.Is(err, ErrCanceled) {
		t.Errorf("неверный тип ошибки: ожидали ErrTimeout или ErrCanceled (в зависимости от того, что произошло раньше), получили %v (тип: %T). Функция должна возвращать одну из этих ошибок при отмене контекста во время выполнения", err, err)
	}
}

func TestWorkTimeoutTiming(t *testing.T) {
	t.Parallel()
	start := time.Now()
	err := Work(context.Background())
	duration := time.Since(start)

	if err == nil {
		t.Fatal("ожидалась ошибка таймаута: функция Work() должна возвращать ошибку ErrTimeout при работе с контекстом без дедлайна. Это указывает на проблему с механизмом таймаута")
	}

	if duration < 90*time.Millisecond || duration > 150*time.Millisecond {
		t.Errorf("неверное время таймаута: выполнение заняло %v, ожидалось примерно 100мс. Таймаут должен происходить примерно через 100мс после начала выполнения. Это указывает на проблему с точностью таймера или механизмом таймаута", duration)
	}

	if !errors.Is(err, ErrTimeout) {
		t.Errorf("неверный тип ошибки: ожидали ErrTimeout, получили %v (тип: %T). При таймауте функция должна возвращать ошибку таймаута, а не другую ошибку", err, err)
	}
}

func TestWorkContextWithDeadline(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(50*time.Millisecond))
	defer cancel()

	err := Work(ctx)
	if err == nil {
		t.Fatal("ожидалась ошибка: функция Work() должна возвращать ошибку при работе с контекстом с дедлайном (50мс). Должна быть возвращена ошибка таймаута или превышения дедлайна. Это указывает на проблему с обработкой дедлайна контекста")
	}
}

func TestWorkMultipleCalls(t *testing.T) {
	t.Parallel()
	const numCalls = 10
	for i := 0; i < numCalls; i++ {
		err := Work(context.Background())
		if err == nil {
			t.Errorf("вызов %d: ожидалась ошибка таймаута, но ошибка не была возвращена. При каждом вызове Work() с контекстом без дедлайна должна возвращаться ошибка ErrTimeout. Это указывает на нестабильность механизма таймаута", i)
		}
		if !errors.Is(err, ErrTimeout) {
			t.Errorf("вызов %d: неверный тип ошибки, ожидали ErrTimeout, получили %v (тип: %T). При множественных вызовах функция должна последовательно возвращать ошибку таймаута", i, err, err)
		}
	}
}

func TestWorkConcurrent(t *testing.T) {
	t.Parallel()
	const numGoroutines = 10
	errs := make([]error, numGoroutines)
	done := make(chan struct{})

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			errs[idx] = Work(context.Background())
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	for i, err := range errs {
		if err == nil {
			t.Errorf("горутина %d: ожидалась ошибка таймаута, но ошибка не была возвращена. При параллельном выполнении Work() должна возвращаться ошибка ErrTimeout. Это указывает на проблему с обработкой ошибок при параллельном выполнении", i)
		}
		if !errors.Is(err, ErrTimeout) {
			t.Errorf("горутина %d: неверный тип ошибки, ожидали ErrTimeout, получили %v (тип: %T). При параллельном выполнении все горутины должны получать одинаковый тип ошибки", i, err, err)
		}
	}
}

func TestWorkCanceledImmediately(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := Work(ctx)
	duration := time.Since(start)

	if err == nil {
		t.Fatal("ожидалась ошибка отмены: функция Work() должна возвращать ошибку ErrCanceled при работе с уже отменённым контекстом. При немедленной отмене контекста функция должна быстро вернуть ошибку отмены. Это указывает на проблему с обработкой уже отменённого контекста")
	}
	if !errors.Is(err, ErrCanceled) {
		t.Errorf("неверный тип ошибки: ожидали ErrCanceled для немедленно отменённого контекста, получили %v (тип: %T). При отмене до начала выполнения функция должна возвращать ошибку отмены", err, err)
	}

	if duration > 50*time.Millisecond {
		t.Errorf("функция должна возвращаться немедленно при отменённом контексте: выполнение заняло %v, ожидалось менее 50мс. При уже отменённом контексте функция не должна ожидать таймаут. Это указывает на проблему с проверкой состояния контекста", duration)
	}
}

func BenchmarkWork(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Work(ctx)
	}
}
