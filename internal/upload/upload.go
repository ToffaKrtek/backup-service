package upload

import (
	"sync"

	"github.com/ToffaKrtek/backup-service/internal/config"
)

func Upload(wg *sync.WaitGroup, item config.S3Item) {
	wg.Add(1)
	defer wg.Done()
	conf := config.Config
}
