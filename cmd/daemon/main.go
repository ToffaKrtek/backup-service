package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ToffaKrtek/backup-service/internal/archive"
	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/ToffaKrtek/backup-service/internal/database"
	"github.com/ToffaKrtek/backup-service/internal/socket"
	"github.com/ToffaKrtek/backup-service/internal/upload"
)

// var jobQueue func()
var mu sync.Mutex
var cancelChan chan struct{}
var jobRunning bool

func main() {
	config.LoadConfig()
	log.Println("Конфигурация загружена:", config.Config)

	go socket.SocketStart(config.UpdateConfigHandler, rerunJobQueueHandler)

	scheduleJob()
	select {}
}

func scheduleJob() {
	fmt.Println("Зашел scheduleJob")
	mu.Lock()
	defer mu.Unlock()

	jobRunning = true
	cancelChan = make(chan struct{}, 1)

	fmt.Println("append jobQueue")
	go func() {
		fmt.Println("Запуск очереди")
		select {
		case <-cancelChan:
			fmt.Println("Джоба отменена")
			jobRunning = false
			return
		default:
			config.LoadConfig()
			now := time.Now()
			startTime := config.Config.GetStartTime()

			if now.Before(startTime) {
				time.Sleep(startTime.Sub(now))
			}
			nextDay := config.Config.StartTime.Add(24 * time.Hour)
			var wg sync.WaitGroup
			files := make(chan config.S3Item)
			for i, schedule := range config.Config.Schedules {
				fmt.Println("Джоба", schedule.ScheduleName)
				if schedule.StartTime.Before(config.Config.StartTime) {
					fmt.Println("Запуск джобы", schedule.ScheduleName)
					if len(schedule.Directories) > 0 {
						wg.Add(1)
						go func() {
							defer wg.Done()
							archive.Archive(&wg, schedule.Directories, files)
						}()
					}
					if len(schedule.DataBases) > 0 {
						wg.Add(1)
						go func() {
							defer wg.Done()
							database.Dump(&wg, schedule.DataBases, files)
						}()
					}
					if schedule.EveryDay {
						config.Config.Schedules[i].StartTime = nextDay
					} else {
						config.Config.Schedules[i].StartTime = config.Config.StartTime.Add(7 * 24 * time.Hour)
					}
				}
				continue
			}
			fmt.Println("Ожидание")
			go func() {
				wg.Wait()
				close(files)
			}()
			for item := range files {
				upload.Upload(&wg, item)
			}
			wg.Wait()
			config.Config.StartTime = nextDay
			config.SaveConfig(false)
			scheduleJob()
		}
	}()
}

func rerunJobQueueHandler(conn net.Conn) {

	fmt.Println("Перезапуск очереди")
	if jobRunning {
		fmt.Println("jobRunning")
		cancelChan <- struct{}{}
		fmt.Println("time.Sleep")
		time.Sleep(100 * time.Millisecond)
	}
	defer conn.Close()
	scheduleJob()
}
