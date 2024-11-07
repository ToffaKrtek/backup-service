package archive

import (
	"log"
	"sync"
	"testing"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/stretchr/testify/assert"
)

var mockConfig = &config.ConfigType{
	Directories: []config.DirectoryConfigType{
		{
			Bucket:  "test-bucket",
			Path:    "test/path/dir1",
			Dirname: "dir1",
		},
		{
			Bucket:  "test-bucket",
			Path:    "test/path/dir2",
			Dirname: "dir2",
		},
	},
}

func TestArchive(t *testing.T) {
	config.Config = mockConfig

	files := make(chan config.S3Item, len(mockConfig.Directories))

	var wg sync.WaitGroup

	Archive(&wg, files)
	wg.Wait()
	close(files)

	expectedItems := map[string]struct{}{
		"dir1": {},
		"dir2": {},
	}

	for item := range files {
		log.Println(item.FilePath)
		assert.Contains(t, expectedItems, item.ObjectName, "Неверный объект в результатах: %s", item.ObjectName)
		delete(expectedItems, item.ObjectName)
	}
	assert.Empty(t, expectedItems, "Не возвращены значения для: %v", expectedItems)
}
