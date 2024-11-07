package upload

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMinioClient struct {
	mock.Mock
}

var mockConfig = &config.ConfigType{
	S3: config.S3ConfigType{
		Endpoint:        "http://localhost:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	},
	ServerName: "test",
}

func (m *MockMinioClient) PutObject(
	ctx context.Context,
	bucket string,
	object string,
	reader io.Reader,
	size int64,
	opts minio.PutObjectOptions,
) (info minio.UploadInfo, err error) {
	args := m.Called(ctx, bucket, object, reader, size, opts)
	log.Println("Отправка объекта")
	return minio.UploadInfo{Size: args.Get(0).(int64)}, args.Error(1)
}

func TestUpload(t *testing.T) {
	var wg sync.WaitGroup

	mockClient := new(MockMinioClient)
	config.Config = mockConfig

	tempFile, err := os.CreateTemp("", "testFile")
	assert.NoError(t, err)
	if err == nil {
		defer os.Remove(tempFile.Name())

		testStr := "test data"
		_, err = tempFile.WriteString(testStr)
		assert.NoError(t, err)
		tempFile.Close()
		item := config.S3Item{
			FilePath:   tempFile.Name(),
			ObjectName: "testobject",
			Bucket:     "testbucket",
		}

		mockClient.On(
			"PutObject",
			mock.Anything,
			item.Bucket,
			mock.Anything,
			mock.Anything,
			int64(len(testStr)),
			mock.Anything,
		).Return(int64(len(testStr)), nil)

		minioClient = mockClient

		Upload(&wg, item)
		wg.Wait()

		mockClient.AssertExpectations(t)
	}
}
