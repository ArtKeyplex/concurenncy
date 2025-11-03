package scheduler

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestEvery(t *testing.T) {
	t.Parallel()
	var count int32
	stop := Every(50*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	time.Sleep(170 * time.Millisecond)
	stop()
	if c := atomic.LoadInt32(&count); c < 3 || c > 4 {
		t.Fatalf("неверное количество выполнений: ожидали 3 или 4 выполнения за 170мс с интервалом 50мс, получили %d. При интервале 50мс за 170мс должно быть примерно 3-4 выполнения (170/50). Это указывает на проблему с таймером или механизмом выполнения", c)
	}
}

func TestEveryStop(t *testing.T) {
	t.Parallel()
	var count int32
	stop := Every(10*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	time.Sleep(25 * time.Millisecond)
	c := atomic.LoadInt32(&count)
	stop()
	time.Sleep(30 * time.Millisecond)
	if atomic.LoadInt32(&count) != c {
		t.Fatal("функция должна быть остановлена")
	}
}

func TestEveryTiming(t *testing.T) {
	t.Parallel()
	var count int32
	interval := 100 * time.Millisecond
	stop := Every(interval, func() { atomic.AddInt32(&count, 1) })

	time.Sleep(250 * time.Millisecond)
	stop()

	c := atomic.LoadInt32(&count)
	if c < 2 || c > 3 {
		t.Fatalf("неверное количество выполнений: ожидали 2-3 выполнения за 250мс с интервалом 100мс, получили %d. При интервале 100мс за 250мс должно быть примерно 2-3 выполнения (250/100). Это указывает на проблему с точностью таймера или механизмом выполнения", c)
	}
}

func TestEveryStopMultipleTimes(t *testing.T) {
	t.Parallel()
	var count int32
	stop := Every(10*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	time.Sleep(20 * time.Millisecond)
	stop()

	time.Sleep(15 * time.Millisecond)
	c1 := atomic.LoadInt32(&count)

	stop()
	stop()

	time.Sleep(30 * time.Millisecond)
	c2 := atomic.LoadInt32(&count)
	if c1 != c2 {
		t.Fatalf("нарушение остановки: счётчик изменился после множественных вызовов stop(): было %d, стало %d. После остановки планировщика функция не должна выполняться, даже при множественных вызовах stop(). Это указывает на проблему с механизмом остановки или утечкой горутины", c1, c2)
	}
}

func TestEveryNoExecutionBeforeInterval(t *testing.T) {
	t.Parallel()
	var count int32
	stop := Every(100*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	time.Sleep(50 * time.Millisecond)
	stop()

	c := atomic.LoadInt32(&count)
	if c != 0 {
		t.Fatalf("преждевременное выполнение: ожидали 0 выполнений до истечения интервала (50мс < 100мс), получили %d. Функция не должна выполняться до истечения первого интервала. Это указывает на проблему с таймером или немедленным выполнением", c)
	}
}

func TestEveryImmediateStop(t *testing.T) {
	t.Parallel()
	var count int32
	stop := Every(10*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	stop()
	time.Sleep(30 * time.Millisecond)

	c := atomic.LoadInt32(&count)
	if c != 0 {
		t.Fatalf("выполнение после немедленной остановки: ожидали 0 выполнений после немедленного вызова stop(), получили %d. При остановке до первого выполнения функция не должна запускаться. Это указывает на проблему с механизмом остановки или условием гонки", c)
	}
}

func TestEveryConcurrentExecution(t *testing.T) {
	t.Parallel()
	var count int32
	stop := Every(50*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
		time.Sleep(10 * time.Millisecond)
	})

	time.Sleep(200 * time.Millisecond)
	stop()

	c := atomic.LoadInt32(&count)
	if c < 3 || c > 4 {
		t.Fatalf("неверное количество выполнений при параллельной работе: ожидали 3-4 выполнения за 200мс с интервалом 50мс, получили %d. Даже при параллельном выполнении функции планировщик должен соблюдать интервал. Это указывает на проблему с обработкой параллельных выполнений", c)
	}
}

func TestEveryMultipleInstances(t *testing.T) {
	t.Parallel()
	var count1, count2 int32
	stop1 := Every(50*time.Millisecond, func() { atomic.AddInt32(&count1, 1) })
	stop2 := Every(50*time.Millisecond, func() { atomic.AddInt32(&count2, 1) })

	time.Sleep(150 * time.Millisecond)
	stop1()
	stop2()

	c1 := atomic.LoadInt32(&count1)
	c2 := atomic.LoadInt32(&count2)

	if c1 < 2 || c1 > 3 {
		t.Errorf("неверное количество выполнений для экземпляра 1: ожидали 2-3 выполнения за 150мс с интервалом 50мс, получили %d. Каждый экземпляр планировщика должен работать независимо. Это указывает на проблему с изоляцией экземпляров", c1)
	}
	if c2 < 2 || c2 > 3 {
		t.Errorf("неверное количество выполнений для экземпляра 2: ожидали 2-3 выполнения за 150мс с интервалом 50мс, получили %d. Каждый экземпляр планировщика должен работать независимо. Это указывает на проблему с изоляцией экземпляров", c2)
	}
}

func TestEveryResourceCleanup(t *testing.T) {
	t.Parallel()
	var count int32
	stop := Every(10*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
	time.Sleep(30 * time.Millisecond)
	stop()

	time.Sleep(50 * time.Millisecond)
	c := atomic.LoadInt32(&count)
	time.Sleep(30 * time.Millisecond)

	if atomic.LoadInt32(&count) != c {
		t.Fatal("выполнение после остановки: функция продолжает выполняться после вызова stop(). Планировщик должен полностью останавливать выполнение после остановки. Это указывает на проблему с механизмом остановки или утечкой ресурсов")
	}
}

func BenchmarkEvery(b *testing.B) {
	var count int32
	for i := 0; i < b.N; i++ {
		stop := Every(1*time.Millisecond, func() { atomic.AddInt32(&count, 1) })
		time.Sleep(5 * time.Millisecond)
		stop()
	}
}
