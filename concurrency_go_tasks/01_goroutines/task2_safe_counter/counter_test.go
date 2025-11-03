package counter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	var c Counter
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 10; j++ {
				c.Inc()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if v := c.Value(); v != 1000 {
		t.Fatalf("неверное итоговое значение счётчика: ожидали 1000 (100 горутин × 10 инкрементов каждая), получили %d. Это указывает на проблему с потокобезопасностью или состоянием гонки при параллельных инкрементах", v)
	}
}

func TestCounterInitialValue(t *testing.T) {
	t.Parallel()
	var c Counter
	if v := c.Value(); v != 0 {
		t.Fatalf("новый счётчик должен быть 0, получено %d", v)
	}
}

func TestCounterSingleIncrement(t *testing.T) {
	t.Parallel()
	var c Counter
	c.Inc()
	if v := c.Value(); v != 1 {
		t.Fatalf("неверное значение после одного инкремента: ожидали 1, получили %d. После единственного вызова Inc() значение должно увеличиться на 1. Это указывает на проблему с операцией инкремента", v)
	}
}

func TestCounterMultipleIncrements(t *testing.T) {
	t.Parallel()
	var c Counter
	for i := 0; i < 100; i++ {
		c.Inc()
	}
	if v := c.Value(); v != 100 {
		t.Fatalf("неверное значение после 100 инкрементов: ожидали 100, получили %d. После последовательных вызовов Inc() значение должно увеличиваться на каждый вызов. Это указывает на проблему с операцией инкремента или потерей обновлений", v)
	}
}

func TestCounterConcurrentReads(t *testing.T) {
	t.Parallel()
	var c Counter
	c.Inc()
	c.Inc()
	c.Inc()

	var wg sync.WaitGroup
	const numReaders = 50
	readValues := make([]int, numReaders)

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			readValues[idx] = c.Value()
		}(i)
	}
	wg.Wait()

	for i, val := range readValues {
		if val != 3 {
			t.Errorf("неверное значение при параллельном чтении: читатель %d прочитал %d вместо 3. Все параллельные читатели должны видеть одинаковое значение (3 инкремента). Это указывает на проблему с синхронизацией чтений или состоянием гонки", i, val)
		}
	}
}

func TestCounterConcurrentReadsAndWrites(t *testing.T) {
	t.Parallel()
	var c Counter
	var wg sync.WaitGroup
	const numWriters = 20
	const numReaders = 20
	const incrementsPerWriter = 10

	var readCount int64

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerWriter; j++ {
				c.Inc()
			}
		}()
	}

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				val := c.Value()
				atomic.AddInt64(&readCount, 1)
				if val < 0 {
					t.Errorf("прочитано отрицательное значение: %d. Счётчик не должен возвращать отрицательные значения, даже при параллельных операциях. Это указывает на проблему с целостностью данных или переполнением", val)
				}
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	expectedValue := numWriters * incrementsPerWriter
	if v := c.Value(); v != expectedValue {
		t.Fatalf("неверное итоговое значение при параллельных чтениях и записях: ожидали %d (20 писателей × 10 инкрементов), получили %d. Даже при параллельных чтениях и записях итоговое значение должно быть корректным. Это указывает на проблему с синхронизацией или потерей обновлений", expectedValue, v)
	}

	if atomic.LoadInt64(&readCount) == 0 {
		t.Error("не было выполнено ни одного чтения: параллельные читатели должны были прочитать значение счётчика. Это указывает на проблему с запуском горутин или синхронизацией")
	}
}

func TestCounterRaceCondition(t *testing.T) {
	t.Parallel()
	var c Counter
	var wg sync.WaitGroup
	const numGoroutines = 1000
	const incrementsPerGoroutine = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				c.Inc()
			}
		}()
	}

	wg.Wait()

	expectedValue := numGoroutines * incrementsPerGoroutine
	if v := c.Value(); v != expectedValue {
		t.Fatalf("обнаружена гонка данных: ожидали %d (1000 горутин × 100 инкрементов каждая), получили %d. При параллельных инкрементах все операции должны быть учтены. Это указывает на проблему с потокобезопасностью, мьютексом или потерей обновлений", expectedValue, v)
	}
}

func TestCounterValueConsistency(t *testing.T) {
	t.Parallel()
	var c Counter
	const numIterations = 100

	for i := 0; i < numIterations; i++ {
		c.Inc()
		val := c.Value()
		if val != i+1 {
			t.Fatalf("неконсистентное значение на итерации %d: ожидали %d, получили %d. После каждого инкремента значение должно быть равно количеству выполненных инкрементов. Это указывает на проблему с атомарностью операций или порядком выполнения", i, i+1, val)
		}
	}
}

func TestCounterMultipleInstances(t *testing.T) {
	t.Parallel()
	var c1, c2 Counter

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			c1.Inc()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			c2.Inc()
		}
	}()

	wg.Wait()

	if v1 := c1.Value(); v1 != 50 {
		t.Errorf("неверное значение для счётчика 1: ожидали 50 (50 инкрементов в отдельной горутине), получили %d. Каждый экземпляр счётчика должен работать независимо. Это указывает на проблему с изоляцией экземпляров или общим состоянием", v1)
	}
	if v2 := c2.Value(); v2 != 50 {
		t.Errorf("неверное значение для счётчика 2: ожидали 50 (50 инкрементов в отдельной горутине), получили %d. Каждый экземпляр счётчика должен работать независимо. Это указывает на проблему с изоляцией экземпляров или общим состоянием", v2)
	}
}

func TestCounterPerformance(t *testing.T) {
	t.Parallel()
	var c Counter
	start := time.Now()

	const numIncrements = 1000000
	for i := 0; i < numIncrements; i++ {
		c.Inc()
	}

	duration := time.Since(start)
	if v := c.Value(); v != numIncrements {
		t.Fatalf("неверное итоговое значение после %d инкрементов: ожидали %d, получили %d. Все инкременты должны быть учтены. Это указывает на проблему с операцией инкремента или потерей обновлений", numIncrements, numIncrements, v)
	}

	if duration > 5*time.Second {
		t.Errorf("проблема с производительностью: %d инкрементов заняли %v. Операции инкремента должны выполняться быстро, даже при большом количестве. Это указывает на проблему с блокировками или синхронизацией", numIncrements, duration)
	}

	t.Logf("Performance: %d increments in %v", numIncrements, duration)
}

func BenchmarkCounterInc(b *testing.B) {
	var c Counter
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Inc()
	}
}

func BenchmarkCounterValue(b *testing.B) {
	var c Counter
	for i := 0; i < 1000; i++ {
		c.Inc()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Value()
	}
}

func BenchmarkCounterConcurrent(b *testing.B) {
	var c Counter
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}
