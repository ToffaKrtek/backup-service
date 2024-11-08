package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/ToffaKrtek/backup-service/internal/socket"
)

type S3Item struct {
	ObjectName string
	FilePath   string
	Bucket     string
}

var configFileName = "backup-service.config.json"
var historyFileName = "backup-service.history.json"

type HistoryType struct {
	Uploads []HistoryUploadItem `json:"uploads"`
}
type HistoryUploadItem struct {
	DateTime time.Time `json:"datetime"`
	Size     string    `json:"size"`
	ItemType string    `json:"type"`
	Status   string    `json:"status"`
	FileName string    `json:"filename"`
}
type ConfigType struct {
	StartTime   time.Time             `json:"start_time"`
	ServerName  string                `json:"server_name"`
	Directories []DirectoryConfigType `json:"directories"`
	DataBases   []DataBaseConfigType  `json:"data_bases"`
	S3          S3ConfigType          `json:"s3"`
}

type S3ConfigType struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Send            bool   `json:"send"`
}

type DirectoryConfigType struct {
	Path    string `json:"path"`
	Dirname string `json:"dirname"`
	Bucket  string `json:"s3_bucket"`
	Active  bool   `json:"active"`
}

var DbTypes = []string{"Невыбрано", "Mysql", "PostgreSql"}
var DbTypesMap = map[string]int{
	"Mysql":      1,
	"PostgreSql": 2,
	"":           0,
}

var ActiveMap = map[bool]string{
	true:  "Активно",
	false: "Не активно",
}

type DataBaseConfigType struct {
	User          string `json:"user"`
	Password      string `json:"password"`
	Address       string `json:"address"`
	ContainerName string `json:"container_name"`
	DataBaseName  string `json:"db_name"`
	IsDocker      bool   `json:"is_docker"`
	Bucket        string `json:"s3_bucket"`
	TypeDB        string `json:"type_db"`
	Active        bool   `json:"active"`
}

func (s S3ConfigType) String() string {
	return fmt.Sprintf("S3ConfigType{Endpoint: %s, AccessKeyID: %s, SecretAccessKey: %s}", s.Endpoint, s.AccessKeyID, s.SecretAccessKey)
}

func (d DirectoryConfigType) String() string {
	return fmt.Sprintf("DirectoryConfigType{Path: %s, Bucket: %s}", d.Path, d.Bucket)
}

func (d DataBaseConfigType) String() string {
	return fmt.Sprintf("DataBaseConfigType{User: %s, Address: %s, ContainerName: %s, IsDocker: %t, Bucket: %s}", d.User, d.Address, d.ContainerName, d.IsDocker, d.Bucket)
}

func (c ConfigType) String() string {
	return fmt.Sprintf("ConfigType{StartTime: %s, ServerName: %s, Directories: %v, DataBases: %v, S3: %s}",
		c.StartTime.Format(time.RFC3339), c.ServerName, c.Directories, c.DataBases, c.S3)
}

func (c ConfigType) GetStartTime() time.Time {
	return c.StartTime
}

var History *HistoryType

func LoadHistory() {

	if _, err := os.Stat(historyFileName); os.IsNotExist(err) {
		emptyHistory := HistoryType{}
		data, err := json.MarshalIndent(emptyHistory, "", "  ")
		if err != nil {
			log.Println("Ошибка сериализации истории отправки:", err)
			return
		}

		if err := os.WriteFile(historyFileName, data, 0644); err != nil {
			log.Println("Ошибка записи истории отправки:", err)
			return
		}

		History = &emptyHistory
		log.Println("История отправки создана:", Config)
		return
	}
	data, err := os.ReadFile(historyFileName)
	if err != nil {
		log.Println("Ошибка чтения истории отправки:", err)
		return
	}
	if err := json.Unmarshal(data, &History); err != nil {
		log.Println("Ошибка парсинга истории отправки:", err)
	}
}

//TODO:: add to history

var Config *ConfigType

func LoadConfig() {
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		log.Println("Файл конфигурации не найден. Создание дефолтной конфигурации.")
		defaultConfig := ConfigType{
			StartTime:   time.Now().Add(24 * time.Hour),
			ServerName:  "default_server",
			Directories: []DirectoryConfigType{},
			DataBases:   []DataBaseConfigType{},
			S3: S3ConfigType{
				Endpoint:        "default_endpoint",
				AccessKeyID:     "default_access_key",
				SecretAccessKey: "default_secret_key",
			},
		}

		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			log.Println("Ошибка сериализации дефолтной конфигурации:", err)
			return
		}

		if err := os.WriteFile(configFileName, data, 0644); err != nil {
			log.Println("Ошибка записи дефолтной конфигурации:", err)
			return
		}

		Config = &defaultConfig
		log.Println("Дефолтная конфигурация создана:", Config)
		return
	}

	data, err := os.ReadFile(configFileName)
	if err != nil {
		log.Println("Ошибка чтения конфигурации:", err)
		return
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		log.Println("Ошибка парсинга конфигурации:", err)
	}
}

func SaveConfig(trigger bool) {
	data, err := json.Marshal(Config)
	if err != nil {
		log.Println("Ошибка сериализации конфигурации:", err)
	}
	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		log.Println("Ошибка записи конфигурации:", err)
	}
	if trigger {
		socket.TriggerSocket()
	}
}

func UpdateConfigHandler(conn net.Conn) {
	defer conn.Close()

	// var newConfig *ConfigType
	// if err := json.NewDecoder(conn).Decode(newConfig); err != nil {
	// 	log.Println("Ошибка декодирования новой конфигурации:", err)
	// 	return
	// }
	// Config = newConfig
	// SaveConfig()
	LoadConfig()
	log.Println("Конфигурация обновлена:", *Config)
}

func GetFileName(template string) string {
	tempDir := os.TempDir()
	timestamp := time.Now().Format("20060102_150405")
	return filepath.Join(tempDir, fmt.Sprintf(template, timestamp))
}
