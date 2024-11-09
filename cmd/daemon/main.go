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
			now := time.Now()
			startTime := config.Config.GetStartTime()

			if now.Before(startTime) {
				time.Sleep(startTime.Sub(now))
			}
			fmt.Println("Запуск джоб")
			var wg sync.WaitGroup
			files := make(chan config.S3Item)
			wg.Add(2)
			go func() {
				defer wg.Done()
				archive.Archive(&wg, files)
			}()
			go func() {
				defer wg.Done()
				database.Dump(&wg, files)
			}()

			fmt.Println("Ожидание")
			go func() {
				wg.Wait()
				close(files)
			}()
			for item := range files {
				upload.Upload(&wg, item)
			}
			wg.Wait()
			config.Config.StartTime = config.Config.StartTime.Add(24 * time.Hour)
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
