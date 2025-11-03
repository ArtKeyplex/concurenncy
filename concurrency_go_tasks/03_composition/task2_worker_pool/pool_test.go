package pool

import (
	"sync"
	"testing"
	"time"
)

func TestRunPool(t *testing.T) {
	t.Parallel()
	jobs := []int{1, 2, 3, 4, 5, 6}
	start := time.Now()
	sum := RunPool(jobs, 3)
	duration := time.Since(start)
	if sum != 21 {
		t.Fatalf("expected 21, got %d", sum)
	}
	if duration >= time.Duration(len(jobs))*10*time.Millisecond {
		t.Fatal("expected parallel execution")
	}
}

func TestRunPoolSingleWorker(t *testing.T) {
	t.Parallel()
	jobs := []int{1, 2, 3}
	if RunPool(jobs, 1) != 6 {
		t.Fatal("неверная сумма при одном воркере")
	}
}

func TestRunPoolEmptyJobs(t *testing.T) {
	t.Parallel()
	jobs := []int{}
	if RunPool(jobs, 3) != 0 {
		t.Fatal("expected 0 for empty jobs")
	}
}

func TestRunPoolZeroWorkers(t *testing.T) {
	t.Parallel()
	jobs := []int{1, 2, 3}
	sum := RunPool(jobs, 0)
	if sum != 6 {
		t.Fatalf("expected 6 with 0 workers (should default to 1), got %d", sum)
	}
}

func TestRunPoolNegativeWorkers(t *testing.T) {
	t.Parallel()
	jobs := []int{1, 2, 3}
	sum := RunPool(jobs, -1)
	if sum != 6 {
		t.Fatalf("expected 6 with negative workers (should default to 1), got %d", sum)
	}
}

func TestRunPoolMoreWorkersThanJobs(t *testing.T) {
	t.Parallel()
	jobs := []int{1, 2}
	sum := RunPool(jobs, 10)
	if sum != 3 {
		t.Fatalf("expected 3, got %d", sum)
	}
}

func TestRunPoolParallelExecution(t *testing.T) {
	t.Parallel()
	jobs := make([]int, 100)
	for i := range jobs {
		jobs[i] = i + 1
	}

	start := time.Now()
	sum := RunPool(jobs, 10)
	duration := time.Since(start)

	expectedSum := 0
	for _, job := range jobs {
		expectedSum += job
	}

	if sum != expectedSum {
		t.Fatalf("expected sum %d, got %d", expectedSum, sum)
	}

	if duration > time.Duration(len(jobs))*time.Millisecond {
		t.Logf("parallel execution took %v (should be fast)", duration)
	}
}

func TestRunPoolConcurrentCalls(t *testing.T) {
	t.Parallel()
	const numCalls = 10
	var wg sync.WaitGroup
	results := make([]int, numCalls)

	jobs := []int{1, 2, 3}
	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = RunPool(jobs, 3)
		}(i)
	}

	wg.Wait()

	expected := 6
	for i, res := range results {
		if res != expected {
			t.Errorf("call %d: expected %d, got %d", i, expected, res)
		}
	}
}

func TestRunPoolLargeJobs(t *testing.T) {
	t.Parallel()
	jobs := make([]int, 1000)
	for i := range jobs {
		jobs[i] = i + 1
	}

	start := time.Now()
	sum := RunPool(jobs, 20)
	duration := time.Since(start)

	expectedSum := 0
	for _, job := range jobs {
		expectedSum += job
	}

	if sum != expectedSum {
		t.Fatalf("expected sum %d, got %d", expectedSum, sum)
	}

	if duration > 5*time.Second {
		t.Errorf("pool took too long: %v", duration)
	}
}

func TestRunPoolSingleJob(t *testing.T) {
	t.Parallel()
	jobs := []int{42}
	sum := RunPool(jobs, 5)
	if sum != 42 {
		t.Fatalf("expected 42, got %d", sum)
	}
}

func TestRunPoolWorkerDistribution(t *testing.T) {
	t.Parallel()
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sum := RunPool(jobs, 3)
	expectedSum := 55
	if sum != expectedSum {
		t.Fatalf("expected %d, got %d", expectedSum, sum)
	}
}

func TestRunPoolCompletes(t *testing.T) {
	t.Parallel()
	jobs := make([]int, 100)
	for i := range jobs {
		jobs[i] = i + 1
	}

	done := make(chan struct{})
	go func() {
		_ = RunPool(jobs, 10)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("pool did not complete (possible deadlock)")
	}
}

func BenchmarkRunPool(b *testing.B) {
	jobs := make([]int, 100)
	for i := range jobs {
		jobs[i] = i + 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RunPool(jobs, 10)
	}
}

func BenchmarkRunPoolSequential(b *testing.B) {
	jobs := make([]int, 100)
	for i := range jobs {
		jobs[i] = i + 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RunPool(jobs, 1)
	}
}
