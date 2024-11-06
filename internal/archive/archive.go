package archive

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/ToffaKrtek/backup-service/internal/config"
)

func Archive(wg *sync.WaitGroup, files chan config.S3Item) {
	conf := config.Config
	for _, dir := range conf.Directories {
		wg.Add(1)
		go func(dir config.DirectoryConfigType) {
			defer wg.Done()
			files <- config.S3Item{
				Bucket:     dir.Bucket,
				FilePath:   zipDirectory(dir.Path, dir.Dirname),
				ObjectName: dir.Dirname,
			}
		}(dir)
	}
}

func zipDirectory(source string, dirname string) string {
	target := config.GetFileName(dirname + "_%s.zip")
	zipFile, err := os.Create(target)
	if err != nil {
		return ""
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(source, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, file)
		if err != nil {
			return err
		}

		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileReader, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		_, err = io.Copy(zipEntry, fileReader)
		return err
	})
	if err != nil {
		return ""
	}

	return target
}
