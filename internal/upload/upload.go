package upload

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Upload(wg *sync.WaitGroup, item config.S3Item) {
	wg.Add(1)
	go func(item config.S3Item) {
		defer wg.Done()
		uploadToMinio(
			item.FilePath,
			item.ObjectName,
			item.Bucket,
		)
	}(item)
	//conf := config.Config
}

func uploadToMinio(
	filePath string,
	objectName string,
	bucket string,
) error {
	conf := config.Config
	minioClient, err := minio.New(conf.S3.Endpoint, &minio.Options{
		// тут скорее всего надо третьим аргументом токен -- TODO::проверить
		Creds:  credentials.NewStaticV4(conf.S3.AccessKeyID, conf.S3.SecretAccessKey, ""),
		Secure: true, //Мб в config
	})
	if err != nil {
		return fmt.Errorf("Не удалось создать клиент: %v", err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Не удалось открыть файл: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Не удалось получить информацию о файле: %v", err)
	}
	curObjName := config.GetFileName(conf.ServerName + objectName + "_%s")
	options := minio.PutObjectOptions{}
	_, err = minioClient.PutObject(
		context.Background(),
		bucket,
		curObjName,
		file,
		fileInfo.Size(),
		options,
	)
	if err != nil {
		return fmt.Errorf("Не удалось отправить файл: %v", err)
	}

	return nil
}
