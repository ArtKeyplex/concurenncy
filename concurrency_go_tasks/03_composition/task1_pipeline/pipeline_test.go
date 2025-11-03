package pipeline

import (
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Parallel()
	nums := []int{1, 2, 3, 4, 5}
	result := Run(nums)
	if result != 110 {
		t.Fatalf("expected 110, got %d", result)
	}
}

func TestRunEmpty(t *testing.T) {
	t.Parallel()
	if Run(nil) != 0 {
		t.Fatal("при пустом списке сумма должна быть 0")
	}
}

func TestRunSingleElement(t *testing.T) {
	t.Parallel()
	nums := []int{5}
	result := Run(nums)
	if result != 50 {
		t.Fatalf("expected 50, got %d", result)
	}
}

func TestRunTwoElements(t *testing.T) {
	t.Parallel()
	nums := []int{2, 3}
	result := Run(nums)
	if result != 26 {
		t.Fatalf("expected 26, got %d", result)
	}
}

func TestRunZeroValues(t *testing.T) {
	t.Parallel()
	nums := []int{0, 0, 0}
	result := Run(nums)
	if result != 0 {
		t.Fatalf("expected 0, got %d", result)
	}
}

func TestRunNegativeValues(t *testing.T) {
	t.Parallel()
	nums := []int{-2, -3}
	result := Run(nums)
	if result != 26 {
		t.Fatalf("expected 26, got %d", result)
	}
}

func TestRunLargeArray(t *testing.T) {
	t.Parallel()
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = i + 1
	}
	result := Run(nums)
	if result <= 0 {
		t.Fatalf("expected positive result, got %d", result)
	}
}

func TestRunPipelineOrder(t *testing.T) {
	t.Parallel()
	nums := []int{1, 2, 3}
	result := Run(nums)
	expected := 28
	if result != expected {
		t.Fatalf("expected %d, got %d", expected, result)
	}
}

func TestRunConcurrentCalls(t *testing.T) {
	t.Parallel()
	const numCalls = 10
	var wg sync.WaitGroup
	results := make([]int, numCalls)

	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			nums := []int{1, 2, 3}
			results[idx] = Run(nums)
		}(i)
	}

	wg.Wait()

	expected := 28
	for i, res := range results {
		if res != expected {
			t.Errorf("call %d: expected %d, got %d", i, expected, res)
		}
	}
}

func TestRunCompletesInTime(t *testing.T) {
	t.Parallel()
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}

	start := time.Now()
	result := Run(nums)
	duration := time.Since(start)

	if result <= 0 {
		t.Errorf("expected positive result, got %d", result)
	}
	if duration > 5*time.Second {
		t.Errorf("pipeline took too long: %v", duration)
	}
}

func TestRunChannelClosure(t *testing.T) {
	t.Parallel()
	nums := []int{1, 2, 3}
	result := Run(nums)
	if result != 28 {
		t.Errorf("expected 28, got %d", result)
	}
}

func BenchmarkRun(b *testing.B) {
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = i + 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Run(nums)
	}
}
