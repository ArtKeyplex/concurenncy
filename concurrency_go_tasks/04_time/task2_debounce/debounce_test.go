package debounce

import (
	"sync"
	"testing"
	"time"
)

func TestDebounce(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(50*time.Millisecond, in)

	go func() {
		in <- 1
		time.Sleep(10 * time.Millisecond)
		in <- 2
		time.Sleep(10 * time.Millisecond)
		in <- 3
		close(in)
	}()

	v, ok := <-out
	if !ok || v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}

	_, ok = <-out
	if ok {
		t.Fatal("channel should be closed")
	}
}

func TestDebounceMultiple(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(20*time.Millisecond, in)
	go func() {
		in <- 1
		time.Sleep(25 * time.Millisecond)
		in <- 2
		close(in)
	}()
	first, ok := <-out
	if !ok || first != 1 {
		t.Fatalf("ожидалось первое значение 1, получено %d", first)
	}
	second, ok := <-out
	if !ok || second != 2 {
		t.Fatalf("ожидалось второе значение 2, получено %d", second)
	}
}

func TestDebounceSingleValue(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(10*time.Millisecond, in)

	go func() {
		in <- 42
		close(in)
	}()

	time.Sleep(20 * time.Millisecond)
	v, ok := <-out
	if !ok || v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestDebounceEmptyChannel(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(10*time.Millisecond, in)

	go func() {
		close(in)
	}()

	time.Sleep(30 * time.Millisecond)
	_, ok := <-out
	if ok {
		t.Fatal("should not receive from empty channel")
	}
}

func TestDebounceRapidValues(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(50*time.Millisecond, in)

	go func() {
		for i := 1; i <= 10; i++ {
			in <- i
			time.Sleep(5 * time.Millisecond)
		}
		close(in)
	}()

	v, ok := <-out
	if !ok || v != 10 {
		t.Fatalf("expected 10 (last value), got %d", v)
	}
}

func TestDebounceMultipleWindows(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(30*time.Millisecond, in)

	go func() {
		in <- 1
		time.Sleep(40 * time.Millisecond)
		in <- 2
		time.Sleep(40 * time.Millisecond)
		in <- 3
		close(in)
	}()

	first, ok := <-out
	if !ok || first != 1 {
		t.Fatalf("expected 1, got %d", first)
	}

	second, ok := <-out
	if !ok || second != 2 {
		t.Fatalf("expected 2, got %d", second)
	}

	third, ok := <-out
	if !ok || third != 3 {
		t.Fatalf("expected 3, got %d", third)
	}
}

func TestDebounceTiming(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(100*time.Millisecond, in)

	start := time.Now()
	go func() {
		in <- 1
		time.Sleep(50 * time.Millisecond)
		in <- 2
		close(in)
	}()

	v, ok := <-out
	duration := time.Since(start)

	if !ok || v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}

	if duration < 100*time.Millisecond {
		t.Errorf("debounce timing incorrect: took %v, expected at least 100ms", duration)
	}
}

func TestDebounceConcurrentInput(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(20*time.Millisecond, in)

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(val int) {
			defer wg.Done()
			in <- val
		}(i)
	}

	go func() {
		wg.Wait()
		close(in)
	}()

	time.Sleep(50 * time.Millisecond)
	count := 0
	for range out {
		count++
	}

	if count != 1 {
		t.Fatalf("expected 1 value, got %d", count)
	}
}

func TestDebounceZeroDuration(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(0, in)

	go func() {
		in <- 1
		close(in)
	}()

	v, ok := <-out
	if !ok || v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}
}

func TestDebounceLargeValues(t *testing.T) {
	t.Parallel()
	in := make(chan int)
	out := Debounce(10*time.Millisecond, in)

	go func() {
		in <- 1000
		in <- 2000
		in <- 3000
		close(in)
	}()

	v, ok := <-out
	if !ok || v != 3000 {
		t.Fatalf("expected 3000, got %d", v)
	}
}

func BenchmarkDebounce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		in := make(chan int)
		out := Debounce(10*time.Millisecond, in)

		go func() {
			for j := 0; j < 10; j++ {
				in <- j
			}
			close(in)
		}()

		for range out {
		}
	}
}
