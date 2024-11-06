package daemon

import (
	"testing"
	"time"
)

func TestRunDaemon(t *testing.T) {
	runJobChan := make(chan struct{})
	restartChan := make(chan struct{})
	now := time.Now()
	startTime := now.Add(3 * time.Second)

	getStartTime := func() time.Time {
		return startTime
	}

	go RunDaemon(getStartTime, runJobChan, restartChan)

	select {
	case <-runJobChan:
		// Сигнал получен, тест пройден
	case <-time.After(5 * time.Second):
		// Если сигнал не получен в течение 5 секунд, тест не пройден
		t.Fatal("Expected signal in runJobChan was not received")
	}

	close(restartChan)
}

func TestRunDaemonRestart(t *testing.T) {
	runJobChan := make(chan struct{})
	restartChan := make(chan struct{})

	now := time.Now()
	startTime := now.Add(2 * time.Second)

	getStartTime := func() time.Time {
		return startTime
	}

	go RunDaemon(getStartTime, runJobChan, restartChan)

	time.Sleep(1 * time.Second)

	close(restartChan)

	select {
	case <-runJobChan:
		// Если сигнал был получен, тест не пройден
		t.Fatal("Expected no signal in runJobChan after restart")
	case <-time.After(1 * time.Second):
		// Если сигнал не был получен в течение 1 секунды, тест пройден
		// Это ожидаемое поведение
	}
}
