package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type S3Item struct {
	ObjectName string
	FilePath   string
	Bucket     string
}

var configFileName = "backup-service.config.json"

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
}

type DirectoryConfigType struct {
	Path    string `json:"path"`
	Dirname string `json:"dirname"`
	Bucket  string `json:"s3_bucket"`
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

var Config *ConfigType

func LoadConfig() {
	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Println("Ошибка чтения конфигурации:", err)
		return
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		log.Println("Ошибка парсинга конфигурации:", err)
	}
}

func SaveConfig() {
	data, err := json.Marshal(Config)
	if err != nil {
		log.Println("Ошибка сериализации конфигурации:", err)
	}
	if err := ioutil.WriteFile(configFileName, data, 0644); err != nil {
		log.Println("Ошибка записи конфигурации:", err)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	var newConfig *ConfigType
	if err := json.NewDecoder(conn).Decode(newConfig); err != nil {
		log.Println("Ошибка декодирования новой конфигурации:", err)
		return
	}
	Config = newConfig
	SaveConfig()
	log.Println("Конфигурация обновлена:", *Config)
}

func GetFileName(template string) string {
	tempDir := os.TempDir()
	timestamp := time.Now().Format("20060102_150405")
	return filepath.Join(tempDir, fmt.Sprintf(template, timestamp))
}
