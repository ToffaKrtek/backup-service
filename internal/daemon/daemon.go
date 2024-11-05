package daemon

import (
	"time"
)

func RunDaemon(
	getStartTime func() time.Time,
	runJobChan chan struct{},
	restartChan chan struct{},
) {
	for {
		now := time.Now()
		startTime := getStartTime()
		nextRun := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			startTime.Hour(),
			startTime.Minute(),
			startTime.Second(),
			0,
			now.Location(),
		)
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}
		sleepDuration := nextRun.Sub(now)

		select {
		case <-time.After(sleepDuration):
			runJobChan <- struct{}{}
		case <-restartChan:
			return
		}
	}
}
