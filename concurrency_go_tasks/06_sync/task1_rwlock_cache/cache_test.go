package cache

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Parallel()
	c := New()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func() {
			c.Set("key", i)
			wg.Done()
		}()
	}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			c.Get("key")
			wg.Done()
		}()
	}
	wg.Wait()
	if _, ok := c.Get("key"); !ok {
		t.Fatal("expected key to exist")
	}
}

func TestCacheSetGet(t *testing.T) {
	t.Parallel()
	c := New()
	c.Set("foo", 42)
	v, ok := c.Get("foo")
	if !ok || v != 42 {
		t.Fatal("значение должно быть 42")
	}
}

func TestCacheEmpty(t *testing.T) {
	t.Parallel()
	c := New()
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected key to not exist")
	}
}

func TestCacheMultipleKeys(t *testing.T) {
	t.Parallel()
	c := New()
	c.Set("key1", "value1")
	c.Set("key2", 42)
	c.Set("key3", true)

	v1, ok1 := c.Get("key1")
	v2, ok2 := c.Get("key2")
	v3, ok3 := c.Get("key3")

	if !ok1 || v1 != "value1" {
		t.Errorf("key1: expected 'value1', got %v", v1)
	}
	if !ok2 || v2 != 42 {
		t.Errorf("key2: expected 42, got %v", v2)
	}
	if !ok3 || v3 != true {
		t.Errorf("key3: expected true, got %v", v3)
	}
}

func TestCacheOverwrite(t *testing.T) {
	t.Parallel()
	c := New()
	c.Set("key", "old")
	c.Set("key", "new")

	v, ok := c.Get("key")
	if !ok || v != "new" {
		t.Errorf("expected 'new', got %v", v)
	}
}

func TestCacheConcurrentReads(t *testing.T) {
	t.Parallel()
	c := New()
	c.Set("key", "value")

	const numReaders = 100
	var wg sync.WaitGroup
	var readErrors int32

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, ok := c.Get("key")
			if !ok || val != "value" {
				atomic.AddInt32(&readErrors, 1)
			}
		}()
	}

	wg.Wait()

	if readErrors > 0 {
		t.Errorf("expected 0 read errors, got %d", readErrors)
	}
}

func TestCacheConcurrentWrites(t *testing.T) {
	t.Parallel()
	c := New()
	const numWriters = 100

	var wg sync.WaitGroup
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			c.Set("key", val)
		}(i)
	}

	wg.Wait()

	val, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to exist after concurrent writes")
	}

	found := false
	for i := 0; i < numWriters; i++ {
		if val == i {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("value %v not in expected range", val)
	}
}

func TestCacheConcurrentReadsAndWrites(t *testing.T) {
	t.Parallel()
	c := New()
	const numOps = 1000
	var wg sync.WaitGroup

	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			c.Set("key", val)
		}(i)
	}

	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Get("key")
		}()
	}

	wg.Wait()

	_, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to exist")
	}
}

func TestCacheReadLockPerformance(t *testing.T) {
	t.Parallel()
	c := New()
	c.Set("key", "value")

	start := time.Now()
	const numReads = 10000
	for i := 0; i < numReads; i++ {
		c.Get("key")
	}
	duration := time.Since(start)

	if duration > 1*time.Second {
		t.Errorf("reads took too long: %v", duration)
	}
	t.Logf("Performance: %d reads in %v", numReads, duration)
}

func TestCacheMultipleInstances(t *testing.T) {
	t.Parallel()
	c1 := New()
	c2 := New()

	c1.Set("key", "value1")
	c2.Set("key", "value2")

	v1, _ := c1.Get("key")
	v2, _ := c2.Get("key")

	if v1 != "value1" {
		t.Errorf("c1: expected 'value1', got %v", v1)
	}
	if v2 != "value2" {
		t.Errorf("c2: expected 'value2', got %v", v2)
	}
}

func TestCacheNilValue(t *testing.T) {
	t.Parallel()
	c := New()
	c.Set("key", nil)

	v, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}

func BenchmarkCacheSet(b *testing.B) {
	c := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("key", i)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := New()
	c.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Get("key")
	}
}

func BenchmarkCacheConcurrent(b *testing.B) {
	c := New()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set("key", 1)
			c.Get("key")
		}
	})
}
