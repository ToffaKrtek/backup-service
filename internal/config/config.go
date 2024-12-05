package config

import (
	"crypto/sha256"
	"encoding/hex"
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
	ItemType   string
	Size       string
}

var configFileName = "/backup-service.config.json"
var configFileNameTmp = "/tmp/backup-service.config.json"
var isTmp = false

func SetTmp(val bool) {
	isTmp = val
}

type ConfigType struct {
	StartTime   time.Time                      `json:"start_time"`
	ServerName  string                         `json:"server_name"`
	Schedules   []ScheduleConfigType           `json:"schedules"`
	Directories map[string]DirectoryConfigType `json:"directories"`
	DataBases   map[string]DataBaseConfigType  `json:"data_bases"`
	S3          S3ConfigType                   `json:"s3"`
}

type S3ConfigType struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Send            bool   `json:"send"`
}

var ActiveMap = map[bool]string{
	true:  "✓",
	false: "✗",
}

func (s S3ConfigType) String() string {
	return fmt.Sprintf("S3ConfigType{Endpoint: %s, AccessKeyID: %s, SecretAccessKey: %s}", s.Endpoint, s.AccessKeyID, s.SecretAccessKey)
}

func (c ConfigType) String() string {
	return fmt.Sprintf("ConfigType{StartTime: %s, ServerName: %s, Schedules: %s, S3: %s}",
		c.StartTime.Format(time.RFC3339), c.ServerName, c.Schedules, c.S3)
}

func (c ConfigType) GetStartTime() time.Time {
	return c.StartTime
}

var Config *ConfigType

func LoadConfig() {
	confFile := configFileName
	if isTmp {
		confFile = configFileNameTmp
	}
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		log.Println("Файл конфигурации не найден. Создание дефолтной конфигурации.")
		nextDay := time.Now().Add(24 * time.Hour)
		defaultConfig := ConfigType{
			StartTime:   nextDay,
			ServerName:  "default_server",
			Directories: map[string]DirectoryConfigType{},
			DataBases:   map[string]DataBaseConfigType{},
			Schedules: []ScheduleConfigType{
				makeSchedule(),
			},
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

		if err := os.WriteFile(confFile, data, 0644); err != nil {
			log.Println("Ошибка записи дефолтной конфигурации:", err)
			return
		}

		Config = &defaultConfig
		log.Println("Дефолтная конфигурация создана:", Config)
		return
	}

	data, err := os.ReadFile(confFile)
	if err != nil {
		log.Println("Ошибка чтения конфигурации:", err)
		return
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		log.Println("Ошибка парсинга конфигурации:", err)
	}
}

func SaveConfig(trigger bool) {
	confFile := configFileName
	if isTmp {
		confFile = configFileNameTmp
	}
	data, err := json.Marshal(Config)
	if err != nil {
		log.Println("Ошибка сериализации конфигурации:", err)
	}
	if err := os.WriteFile(confFile, data, 0644); err != nil {
		log.Println("Ошибка записи конфигурации:", err)
	}
	if trigger {
		socket.TriggerSocket()
	}
}

func UpdateConfigHandler(conn net.Conn) {
	defer conn.Close()
	LoadConfig()
	log.Println("Конфигурация обновлена:", *Config)
}

func GetFileName(template string) string {
	tempDir := os.TempDir()
	timestamp := time.Now().Format("20060102_150405")
	return filepath.Join(tempDir, fmt.Sprintf(template, timestamp))
}

func GetFileSize(filePath string) string {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return ""
	}
	fileSize := fileInfo.Size()

	gbSize := float64(fileSize) / (1024 * 1024 * 1024)

	if gbSize < 1 {
		return "<1 Гб"
	}
	return fmt.Sprintf("%.2f Гб", gbSize)
}

func getIndex() string {
	now := time.Now().String()
	hash := sha256.New()
	hash.Write([]byte(now))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func UpdateConfig(conf ConfigType, withDay bool) {
	LoadConfig()
	withTrigger := false
	if conf.StartTime.Minute() != Config.StartTime.Minute() ||
		conf.StartTime.Hour() != Config.StartTime.Hour() {
		withTrigger = true
		for i, schedule := range conf.Schedules {
			day := schedule.StartTime.Day()
			if withDay {
				day = conf.StartTime.Day()
			}
			conf.Schedules[i].StartTime = time.Date(
				schedule.StartTime.Year(),
				schedule.StartTime.Month(),
				day,
				conf.StartTime.Hour(),
				conf.StartTime.Minute(),
				conf.StartTime.Second(),
				0,
				conf.StartTime.Location(),
			)
		}
	}
	Config = &conf
	SaveConfig(withTrigger)
}
