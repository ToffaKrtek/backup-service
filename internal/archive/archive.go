package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ToffaKrtek/backup-service/internal/config"
)

func Archive(wg *sync.WaitGroup, dirs map[string]config.DirectoryConfigType, files chan config.S3Item) {
	fmt.Printf("Запуск архивации для %d директорий", len(dirs))
	for _, dir := range dirs {
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

func zipRecentFiles(source string, target string, days int) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("ошибка создания ZIP файла: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	now := time.Now()
	cutoffTime := now.AddDate(0, 0, -days)

	err = filepath.Walk(source, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка при обходе файлов: %w", err)
		}

		// Получаем относительный путь к файлу
		relPath, err := filepath.Rel(source, file)
		if err != nil {
			return fmt.Errorf("ошибка получения относительного пути: %w", err)
		}

		// Проверяем, является ли это директорией или символической ссылкой
		if info.IsDir() || (info.Mode()&os.ModeSymlink) == os.ModeSymlink {
			// Создаем запись для директории
			_, err := zipWriter.Create(relPath + "/")
			if err != nil {
				return fmt.Errorf("ошибка создания записи для директории %s: %w", relPath, err)
			}
			return nil // Возвращаем nil, чтобы продолжить обход
		}

		// Проверяем, был ли файл изменен или создан в последние N дней
		if info.ModTime().After(cutoffTime) {
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
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("ошибка при создании ZIP архива: %w", err)
	}

	return nil
}
