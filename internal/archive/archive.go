package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/ToffaKrtek/backup-service/internal/config"
)

func Archive(wg *sync.WaitGroup, files chan config.S3Item) {
	conf := config.Config
	fmt.Printf("Запуск архивации для %d директорий", len(conf.Directories))
	for _, dir := range conf.Directories {
		if dir.Active {
			wg.Add(1)
			go func(dir config.DirectoryConfigType) {
				defer wg.Done()
				defer fmt.Println("Закончена архивация")
				archivePath := config.GetFileName(dir.Dirname + "_%s.zip")
				if err := zipDirectory(dir.Path, archivePath); err == nil {
					item := config.S3Item{
						Bucket:     dir.Bucket,
						FilePath:   archivePath,
						ObjectName: dir.Dirname + ".zip",
						ItemType:   "Архив",
						Size:       config.GetFileSize(archivePath),
					}
					fmt.Println(item.FilePath)
					files <- item
				} else {
					fmt.Printf("Ошибка: %v", err)
				}
			}(dir)
		}
	}
}
func zipDirectory(source string, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("ошибка создания ZIP файла: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(source, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка при обходе файлов: %w", err)
		}

		// Получаем относительный путь к файлу
		relPath, err := filepath.Rel(source, file)
		if err != nil {
			return fmt.Errorf("ошибка получения относительного пути: %w", err)
		}

		if info.IsDir() || (info.Mode()&os.ModeSymlink) == os.ModeSymlink {
			// Создаем запись для директории
			_, err := zipWriter.Create(relPath + "/")
			if err != nil {
				return fmt.Errorf("ошибка создания записи для директории %s: %w", relPath, err)
			}
			return nil // Возвращаем nil, чтобы продолжить обход
		}

		// Создаем запись для файла
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return fmt.Errorf("ошибка создания записи для файла %s: %w", relPath, err)
		}

		// Открываем файл
		fileReader, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("ошибка открытия файла %s: %w", file, err)
		}
		defer fileReader.Close()

		// Копируем содержимое файла в zip
		_, err = io.Copy(zipEntry, fileReader)
		if err != nil {
			return fmt.Errorf("ошибка копирования содержимого файла %s: %w", file, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("ошибка при создании ZIP архива: %w", err)
	}

	return nil
}
func addFileToZip(zipWriter *zip.Writer, file string, relPath string) error {
	zipEntry, err := zipWriter.Create(relPath)
	if err != nil {
		return fmt.Errorf("ошибка создания записи для файла %s: %w", relPath, err)
	}

	fileReader, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла %s: %w", file, err)
	}
	defer fileReader.Close()

	_, err = io.Copy(zipEntry, fileReader)
	if err != nil {
		return fmt.Errorf("ошибка копирования содержимого файла %s: %w", file, err)
	}
	return nil
}
