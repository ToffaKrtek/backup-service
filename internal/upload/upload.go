package upload

import (
	"context"
	"fmt"
	"io"
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

type MinioClientInterface interface {
	PutObject(
		ctx context.Context,
		bucketName string,
		objectName string,
		reader io.Reader,
		objectSize int64,
		opts minio.PutObjectOptions,
	) (info minio.UploadInfo, err error)
}

var minioClient MinioClientInterface

func newClient(endpoint string, opts *minio.Options) (MinioClientInterface, error) {
	if minioClient == nil {
		return minio.New(endpoint, opts)
	}
	return minioClient, nil
}

func uploadToMinio(
	filePath string,
	objectName string,
	bucket string,
) error {
	conf := config.Config
	minioClient, err := newClient(conf.S3.Endpoint, &minio.Options{
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
