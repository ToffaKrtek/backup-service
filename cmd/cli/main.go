package main

import (
	"fmt"
	"log"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/rivo/tview"
)

func main() {
	config.LoadConfig()

	app := tview.NewApplication()

	form := tview.NewForm().
		//AddInputField("Время запуска", config.Config.StartTime, 20, nil, nil).
		AddInputField("Имя сервера", config.Config.ServerName, 20, nil, nil).
		AddInputField("URL", config.Config.S3.Endpoint, 40, nil, nil).
		AddInputField("AccessKeyID", config.Config.S3.AccessKeyID, 40, nil, nil).
		AddInputField("SecretAccessKey", config.Config.S3.SecretAccessKey, 40, nil, nil).
		AddButton("Сохранить", func() {
			// Логика сохранения конфигурации
			fmt.Println("Конфигурация сохранена.")
		}).
		AddButton("Выход", func() {
			app.Stop()
		}).
		AddButton("Запустить сейчас", func() {
			// Логика запуска сервиса
			fmt.Println("Запуск сервиса...")
		}).
		AddButton("Остановить", func() {
			// Логика остановки сервиса
			fmt.Println("Остановка сервиса...")
		})

	form.SetTitle("Конфигурация бекап сервиса").SetTitleAlign(tview.AlignLeft).SetBorder(true)
	if err := app.SetRoot(form, true).Run(); err != nil {
		log.Fatalf("Ошибка в приложении: %v", err)
	}
}

func runNow() error {
	return nil
}
func killDaemon() error {
	return nil
}
func saveConfig() error {
	return nil
}
