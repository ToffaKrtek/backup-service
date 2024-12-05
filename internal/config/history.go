package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var historyFileName = "/backup-service.history.json"
var historyFileNameTmp = "/tmp/backup-service.history.json"

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

var History *HistoryType

func LoadHistory() {
	histFile := historyFileName
	if isTmp {
		histFile = historyFileNameTmp
	}
	if _, err := os.Stat(histFile); os.IsNotExist(err) {
		emptyHistory := HistoryType{}
		data, err := json.MarshalIndent(emptyHistory, "", "  ")
		if err != nil {
			log.Println("Ошибка сериализации истории отправки:", err)
			return
		}

		if err := os.WriteFile(histFile, data, 0644); err != nil {
			log.Println("Ошибка записи истории отправки:", err)
			return
		}

		History = &emptyHistory
		log.Println("История отправки создана:", Config)
		return
	}
	data, err := os.ReadFile(histFile)
	if err != nil {
		log.Println("Ошибка чтения истории отправки:", err)
		return
	}
	if err := json.Unmarshal(data, &History); err != nil {
		log.Println("Ошибка парсинга истории отправки:", err)
	}
}

// TODO:: add to history
func SaveHistoryItem(item HistoryUploadItem) {
	histFile := historyFileName
	if isTmp {
		histFile = historyFileNameTmp
	}
	if History == nil {
		LoadHistory()
	}
	History.Uploads = append(History.Uploads, item)

	if len(History.Uploads) > 5 {
		History.Uploads = History.Uploads[len(History.Uploads)-5:]
	}

	data, err := json.MarshalIndent(History, "", "  ")
	if err != nil {
		log.Println("Ошибка сериализации истории отправки:", err)
		return
	}

	if err := os.WriteFile(histFile, data, 0644); err != nil {
		log.Println("Ошибка записи истории отправки:", err)
	}
}
