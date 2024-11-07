package main

import (
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

var jobQueue []func()
var mu sync.Mutex

func main() {
	config.LoadConfig()
	log.Println("Конфигурация загружена:", config.Config)

	go socket.SocketStart(config.UpdateConfigHandler, rerunJobQueueHandler)

	scheduleJob()
	select {}
}

func scheduleJob() {
	mu.Lock()
	defer mu.Unlock()

	jobQueue = append(jobQueue, func() {
		now := time.Now()
		startTime := config.Config.GetStartTime()

		if now.Before(startTime) {
			time.Sleep(startTime.Sub(now))
		}

		var wg sync.WaitGroup
		files := make(chan config.S3Item)
		go archive.Archive(&wg, files)
		go database.Dump(&wg, files)

		go func() {
			wg.Wait()
			close(files)
		}()
		for item := range files {
			upload.Upload(&wg, item)
			wg.Wait()
		}
		config.Config.StartTime = config.Config.StartTime.Add(24 * time.Hour)
		config.SaveConfig()
		scheduleJob()
	})
	go jobQueue[len(jobQueue)-1]()
}

func rerunJobQueueHandler(conn net.Conn) {
	defer conn.Close()
	scheduleJob()
}
