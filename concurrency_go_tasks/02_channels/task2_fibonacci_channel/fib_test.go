package fibonacci

import (
	"testing"
	"time"
)

func TestFib(t *testing.T) {
	t.Parallel()
	ch := Fib(6)
	expected := []int{0, 1, 1, 2, 3, 5}
	i := 0
	for v := range ch {
		if v != expected[i] {
			t.Fatalf("expected %d, got %d", expected[i], v)
		}
		i++
	}
	if i != len(expected) {
		t.Fatalf("expected %d numbers, got %d", len(expected), i)
	}
}

func TestFibZero(t *testing.T) {
	t.Parallel()
	ch := Fib(0)
	for range ch {
		t.Fatal("канал должен быть пуст")
	}
}

func TestFibOne(t *testing.T) {
	t.Parallel()
	ch := Fib(1)
	values := make([]int, 0)
	for v := range ch {
		values = append(values, v)
	}
	if len(values) != 1 || values[0] != 0 {
		t.Fatalf("expected [0], got %v", values)
	}
}

func TestFibTwo(t *testing.T) {
	t.Parallel()
	ch := Fib(2)
	expected := []int{0, 1}
	values := make([]int, 0)
	for v := range ch {
		values = append(values, v)
	}
	if len(values) != 2 {
		t.Fatalf("expected 2 numbers, got %d", len(values))
	}
	for i, v := range expected {
		if values[i] != v {
			t.Errorf("index %d: expected %d, got %d", i, v, values[i])
		}
	}
}

func TestFibNegative(t *testing.T) {
	t.Parallel()
	ch := Fib(-5)
	count := 0
	for range ch {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 values for negative input, got %d", count)
	}
}

func TestFibLarge(t *testing.T) {
	t.Parallel()
	n := 20
	ch := Fib(n)
	values := make([]int, 0)
	for v := range ch {
		values = append(values, v)
	}

	if len(values) != n {
		t.Fatalf("expected %d numbers, got %d", n, len(values))
	}

	if values[0] != 0 {
		t.Errorf("first value should be 0, got %d", values[0])
	}
	if n > 1 && values[1] != 1 {
		t.Errorf("second value should be 1, got %d", values[1])
	}
	for i := 2; i < n; i++ {
		expected := values[i-1] + values[i-2]
		if values[i] != expected {
			t.Errorf("index %d: expected %d, got %d (should be sum of previous two)", i, expected, values[i])
		}
	}
}

func TestFibChannelClosure(t *testing.T) {
	t.Parallel()
	ch := Fib(5)
	count := 0
	for range ch {
		count++
	}

	_, ok := <-ch
	if ok {
		t.Fatal("channel should be closed after reading all values")
	}

	if count != 5 {
		t.Errorf("expected 5 values, got %d", count)
	}
}

func TestFibConcurrentReads(t *testing.T) {
	t.Parallel()
	ch := Fib(10)
	values := make([]int, 0)

	done := make(chan struct{})
	go func() {
		for v := range ch {
			values = append(values, v)
		}
		close(done)
	}()

	select {
	case <-done:
		if len(values) != 10 {
			t.Errorf("expected 10 values, got %d", len(values))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("channel reading did not complete")
	}
}

func TestFibCompletesInTime(t *testing.T) {
	t.Parallel()
	start := time.Now()
	ch := Fib(100)
	count := 0
	for range ch {
		count++
	}
	duration := time.Since(start)

	if count != 100 {
		t.Errorf("expected 100 values, got %d", count)
	}
	if duration > 1*time.Second {
		t.Errorf("fibonacci generation took too long: %v", duration)
	}
}

func TestFibMultipleCalls(t *testing.T) {
	t.Parallel()
	const numCalls = 10
	for i := 0; i < numCalls; i++ {
		ch := Fib(5)
		count := 0
		for range ch {
			count++
		}
		if count != 5 {
			t.Errorf("call %d: expected 5 values, got %d", i, count)
		}
	}
}

func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := Fib(100)
		for range ch {
		}
	}
}
