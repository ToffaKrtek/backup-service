package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestLoadConfigAndSaveConfig(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "backup-service.config.json")
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Удаляем файл после теста

	// Создаем тестовую конфигурацию
	testConfig := &ConfigType{
		StartTime:   time.Now(),
		ServerName:  "TestServer",
		Directories: []DirectoryConfigType{{Path: "/test/path", Dirname: "test", Bucket: "test-bucket"}},
		DataBases:   []DataBaseConfigType{{User: "testuser", Password: "testpass", Address: "localhost", ContainerName: "test-container", DataBaseName: "testdb", IsDocker: false, Bucket: "test-bucket", TypeDB: "mysql"}},
		S3:          S3ConfigType{Endpoint: "http://s3.test.com", AccessKeyID: "testAccessKey", SecretAccessKey: "testSecretKey"},
	}

	// Сохраняем тестовую конфигурацию в временный файл
	Config = testConfig
	configFileName = tempFile.Name() // Устанавливаем имя файла конфигурации на временный файл
	SaveConfig()

	// Загружаем конфигурацию из временного файла
	LoadConfig()

	// Проверяем, что загруженная конфигурация совпадает с тестовой
	if Config.ServerName != testConfig.ServerName {
		t.Errorf("Ожидалось %s, получено %s", testConfig.ServerName, Config.ServerName)
	}
	if len(Config.Directories) != len(testConfig.Directories) || Config.Directories[0].Path != testConfig.Directories[0].Path {
		t.Errorf("Ожидалось %v, получено %v", testConfig.Directories, Config.Directories)
	}
	if len(Config.DataBases) != len(testConfig.DataBases) || Config.DataBases[0].User != testConfig.DataBases[0].User {
		t.Errorf("Ожидалось %v, получено %v", testConfig.DataBases, Config.DataBases)
	}
	if Config.S3.Endpoint != testConfig.S3.Endpoint {
		t.Errorf("Ожидалось %s, получено %s", testConfig.S3.Endpoint, Config.S3.Endpoint)
	}
}

func TestGetFileName(t *testing.T) {
	template := "backup_%s.zip"
	fileName := GetFileName(template)

	// Проверяем, что имя файла содержит временную метку
	if time.Now().Format("20060102_150405") != fileName[len(fileName)-19:len(fileName)-4] {
		t.Errorf("Имя файла не соответствует ожидаемому шаблону: %s", fileName)
	}
}
