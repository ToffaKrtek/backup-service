package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/rivo/tview"
)

func main() {
	config.LoadConfig()

	app := tview.NewApplication()

	form := tview.NewForm().
		//AddInputField("Время запуска", config.Config.StartTime, 20, nil, nil).
		AddInputField("Имя сервера", config.Config.ServerName, 20, nil, nil).
		AddInputField(
			"Час запуска",
			strconv.Itoa(config.Config.StartTime.Hour()),
			20,
			func(txt string, l rune) bool {
				num, err := strconv.ParseInt(txt, 10, 0)
				if err != nil {
					return false
				}
				return num <= 23 && num >= 0
			},
			nil,
		).
		AddInputField(
			"Минута запуска",
			strconv.Itoa(config.Config.StartTime.Minute()),
			20,
			func(txt string, l rune) bool {
				num, err := strconv.ParseInt(txt, 10, 0)
				if err != nil {
					return false
				}
				return num <= 59 && num >= 0
			},
			nil,
		).
		AddInputField("URL", config.Config.S3.Endpoint, 40, nil, nil).
		AddInputField("AccessKeyID", config.Config.S3.AccessKeyID, 40, nil, nil).
		AddPasswordField("SecretAccessKey", config.Config.S3.SecretAccessKey, 40, '*', nil).
		AddButton("Сохранить", func() {
			// Логика сохранения конфигурации
			fmt.Println(" Конфигурация сохранена.")
		}).
		AddButton(
			fmt.Sprintf("Директории для архивации (%d)", len(config.Config.Directories)),
			func() {
				// Логика сохранения конфигурации
				fmt.Println(" Конфигурация сохранена.")
			}).
		AddButton(
			fmt.Sprintf("ДБ для создания дампов (%d)", len(config.Config.DataBases)),
			func() {
				// Логика сохранения конфигурации
				fmt.Println(" Конфигурация сохранена.")
			}).
		AddButton("Запустить сейчас", func() {
			// Логика запуска сервиса
			fmt.Println(" Запуск сервиса...")
		}).
		AddButton("Остановить", func() {
			// Логика остановки сервиса
			fmt.Println(" Остановка сервиса...")
		}).
		AddButton("Выход", func() {
			app.Stop()
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
