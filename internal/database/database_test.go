package database

import (
	"io"
	"os"
	"sync"
	"testing"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommand struct {
	mock.Mock
}

func (m *MockCommand) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCommand) Stdout() *os.File {
	return nil
}
func (m *MockCommand) SetStdout(stdout io.Writer) {}
func (m *MockCommand) SetEnv(env []string)        {}

var mockConfig = &config.ConfigType{
	DataBases: []config.DataBaseConfigType{
		{
			TypeDB:        "mysql",
			IsDocker:      true,
			ContainerName: "mysql_container",
			DataBaseName:  "test_db",
			User:          "user",
			Password:      "password",
			Bucket:        "bucket_name",
		},
		{
			TypeDB:        "postgre",
			IsDocker:      true,
			ContainerName: "postgre_container",
			DataBaseName:  "test_db",
			User:          "user",
			Password:      "password",
			Bucket:        "bucket_name",
		},
	},
}

func TestDump(t *testing.T) {

	mockCmd := new(MockCommand)
	mockCmd.On("Run").Return(nil).Times(2)
	execCommand = mockCmd
	files := make(chan config.S3Item, 2)
	var wg sync.WaitGroup

	config.Config = mockConfig
	Dump(&wg, files)
	wg.Wait()
	close(files)
	var items []config.S3Item

	for file := range files {
		items = append(items, file)
	}

	assert.Len(t, items, 2)
	if len(items) > 0 {
		assert.Equal(t, "bucket_name", items[0].Bucket)
		assert.NotEmpty(t, items[0].FilePath)
	}
}
