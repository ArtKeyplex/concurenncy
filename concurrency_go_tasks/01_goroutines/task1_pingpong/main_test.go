package main

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPingPong(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	PingPong(&buf)
	output := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(output) != 10 {
		t.Fatalf("неверное количество строк вывода: ожидали 10 строк (5 ping и 5 pong), получили %d строк", len(output))
	}
	pingCount, pongCount := 0, 0
	for i, line := range output {
		if i%2 == 0 && line != "ping" {
			t.Errorf("неверная последовательность на позиции %d: ожидали 'ping', получили '%s'. Это означает нарушение порядка чередования", i, line)
		}
		if i%2 == 1 && line != "pong" {
			t.Errorf("неверная последовательность на позиции %d: ожидали 'pong', получили '%s'. Это означает нарушение порядка чередования", i, line)
		}
		if line == "ping" {
			pingCount++
		} else if line == "pong" {
			pongCount++
		}
	}
	if pingCount != 5 || pongCount != 5 {
		t.Fatalf("неверное количество сообщений: ожидали 5 ping и 5 pong, получили %d ping и %d pong. Это означает, что функция не выполнилась полностью", pingCount, pongCount)
	}
}

func TestPingPongStartsWithPing(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	PingPong(&buf)
	output := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(output) == 0 || output[0] != "ping" {
		t.Fatal("первой строкой должен быть ping")
	}
}

func TestPingPongCorrectSequence(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	PingPong(&buf)
	output := strings.Split(strings.TrimSpace(buf.String()), "\n")

	expected := []string{"ping", "pong", "ping", "pong", "ping", "pong", "ping", "pong", "ping", "pong"}
	if len(output) != len(expected) {
		t.Fatalf("неверное количество строк: ожидали %d строк (точная последовательность ping-pong), получили %d строк", len(expected), len(output))
	}
	for i, expectedLine := range expected {
		if output[i] != expectedLine {
			t.Errorf("нарушение последовательности на позиции %d: ожидали '%s', получили '%s'. Правильная последовательность должна быть: ping, pong, ping, pong...", i, expectedLine, output[i])
		}
	}
}

func TestPingPongConcurrentExecution(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup
	const numRuns = 10

	for i := 0; i < numRuns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var buf bytes.Buffer
			PingPong(&buf)
			output := strings.Split(strings.TrimSpace(buf.String()), "\n")
			if len(output) != 10 {
				t.Errorf("ошибка при параллельном выполнении: ожидали 10 строк вывода, получили %d. Это может указывать на проблему синхронизации между горутинами", len(output))
			}
		}()
	}

	wg.Wait()
}

func TestPingPongCompletesInTime(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	done := make(chan struct{})

	go func() {
		PingPong(&buf)
		close(done)
	}()

	select {
	case <-done:
		output := strings.Split(strings.TrimSpace(buf.String()), "\n")
		if len(output) != 10 {
			t.Errorf("неверное количество строк после завершения: ожидали 10 строк вывода, получили %d. Функция завершилась, но вывод неполный", len(output))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("функция PingPong не завершилась в течение 5 секунд: возможна взаимная блокировка (deadlock) между горутинами или проблема с синхронизацией каналов")
	}
}

func TestPingPongNoDeadlock(t *testing.T) {
	t.Parallel()
	const numRuns = 50
	var wg sync.WaitGroup

	for i := 0; i < numRuns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var buf bytes.Buffer
			PingPong(&buf)
			output := strings.Split(strings.TrimSpace(buf.String()), "\n")
			if len(output) != 10 {
				t.Errorf("ошибка в запуске %d: ожидали 10 строк вывода, получили %d. Это может указывать на проблему при множественных параллельных вызовах", i, len(output))
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("параллельные выполнения не завершились в течение 10 секунд: возможна взаимная блокировка (deadlock) или проблема с синхронизацией при множественных одновременных вызовах функции PingPong")
	}
}

func TestPingPongNoGoroutineLeak(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer

	start := time.Now()
	PingPong(&buf)
	duration := time.Since(start)

	if duration > 1*time.Second {
		t.Errorf("функция PingPong выполнялась слишком долго (%v): возможна утечка горутин или блокировка. Нормальное выполнение должно занимать миллисекунды", duration)
	}

	output := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(output) != 10 {
		t.Errorf("неверное количество строк вывода: ожидали 10 строк, получили %d. Это может указывать на незавершённые горутины", len(output))
	}
}
