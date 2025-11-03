package pipelinectx

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Parallel()
	result, err := Run(context.Background(), []int{1, 2, 3, 4, 5})
	if err != nil {
		t.Fatalf("не ожидали ошибку при выполнении конвейера: получена ошибка %v", err)
	}
	if result != 30 {
		t.Fatalf("неверный результат суммирования: ожидали сумму 30, получили %d", result)
	}
}

func TestRunCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := Run(ctx, []int{1, 2})
	if err == nil {
		t.Fatal("ожидалась ошибка отмены контекста: контекст был отменён до начала выполнения, но ошибка не возвращена")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("неверный тип ошибки: ожидали context.Canceled, получили %v (тип: %T)", err, err)
	}
	if result != 0 {
		t.Fatalf("неверный результат при отменённом контексте: ожидали 0, получили %d", result)
	}
}

func TestRunEmpty(t *testing.T) {
	t.Parallel()
	res, err := Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("не ожидали ошибку при пустом списке: получена ошибка %v", err)
	}
	if res != 0 {
		t.Fatalf("неверный результат для пустого списка: ожидали 0, получили %d", res)
	}
}

func TestRunSingleElement(t *testing.T) {
	t.Parallel()
	result, err := Run(context.Background(), []int{5})
	if err != nil {
		t.Fatalf("не ожидали ошибку при обработке одного элемента: получена ошибка %v", err)
	}
	if result != 10 {
		t.Fatalf("неверный результат для одного элемента: ожидали 10 (5*2), получили %d", result)
	}
}

func TestRunCancellationMidway(t *testing.T) {
	t.Parallel()

	nums := make([]int, 10000000)
	for i := range nums {
		nums[i] = i + 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	result, err := Run(ctx, nums)

	if err != nil {
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("неверный тип ошибки при отмене контекста: ожидали context.Canceled или context.DeadlineExceeded, получили %v (тип: %T)", err, err)
		}
	}
	if result < 0 {
		t.Errorf("неверный результат при отмене: ожидали неотрицательное значение, получили %d", result)
	}
}

func TestRunTimeout(t *testing.T) {
	t.Parallel()

	nums := make([]int, 100000000)
	for i := range nums {
		nums[i] = i + 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	result, err := Run(ctx, nums)

	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			t.Errorf("неверный тип ошибки при таймауте: ожидали context.DeadlineExceeded или context.Canceled, получили %v (тип: %T)", err, err)
		}
	}
	if result < 0 {
		t.Errorf("неверный результат при таймауте: ожидали неотрицательное значение, получили %d", result)
	}
}

func TestRunCompleteProcessing(t *testing.T) {
	t.Parallel()
	nums := []int{1, 2, 3, 4, 5}
	result, err := Run(context.Background(), nums)
	if err != nil {
		t.Fatalf("не ожидали ошибку при полной обработке: получена ошибка %v", err)
	}
	expected := 30
	if result != expected {
		t.Fatalf("неверный результат полной обработки: ожидали %d, получили %d", expected, result)
	}
}

func TestRunNegativeValues(t *testing.T) {
	t.Parallel()
	nums := []int{-1, -2, -3}
	result, err := Run(context.Background(), nums)
	if err != nil {
		t.Fatalf("не ожидали ошибку при обработке отрицательных значений: получена ошибка %v", err)
	}
	expected := -12
	if result != expected {
		t.Fatalf("неверный результат для отрицательных значений: ожидали %d, получили %d", expected, result)
	}
}

func TestRunAlreadyCanceledContext(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := Run(ctx, []int{1, 2, 3})
	if err == nil {
		t.Fatal("ожидалась ошибка от отменённого контекста: контекст был отменён до вызова Run, но ошибка не возвращена")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("неверный тип ошибки: ожидали context.Canceled, получили %v (тип: %T). Детали ошибки: %v", err, err, err)
	}
	if result != 0 {
		t.Errorf("неверный результат для отменённого контекста: ожидали 0 (обработка не должна была начаться), получили %d. Это означает, что обработка началась несмотря на отменённый контекст", result)
	}
}

func TestRunContextHierarchy(t *testing.T) {
	t.Parallel()
	parentCtx := context.Background()
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	nums := []int{1, 2, 3, 4, 5}
	result, err := Run(ctx, nums)
	if err != nil {
		t.Fatalf("не ожидали ошибку при работе с иерархией контекстов: получена ошибка %v", err)
	}
	if result != 30 {
		t.Fatalf("неверный результат при иерархии контекстов: ожидали 30, получили %d", result)
	}
}

func BenchmarkRun(b *testing.B) {
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = i + 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Run(context.Background(), nums)
	}
}
