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

	pages := tview.NewPages()

	mainForm := tview.NewForm()
	updateMainForm(mainForm, pages, app)

	dirForm := tview.NewForm()
	updateDirForm(dirForm, pages, app)

	dbForm := tview.NewForm()
	updateDBForm(dbForm, pages, app)
	pages.AddPage("main", mainForm, true, true)
	pages.AddPage("dirs", dirForm, true, false)
	pages.AddPage("dbs", dbForm, true, false)
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		log.Fatalf("Ошибка в приложении: %v", err)
	}
}

func updateMainForm(mainForm *tview.Form, pages *tview.Pages, app *tview.Application) {
	mainForm.Clear(true)
	mainForm.AddInputField("Имя сервера", config.Config.ServerName, 20, nil, nil).
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
				pages.SwitchToPage("dirs")
				fmt.Println(" Настройка архивируемых директорий.")
			}).
		AddButton(
			fmt.Sprintf("Дампы ДБ (%d)", len(config.Config.DataBases)),
			func() {
				// Логика сохранения конфигурации
				pages.SwitchToPage("dbs")
				fmt.Println(" Настройка создания дампов ДБ.")
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
	mainForm.SetTitle("Конфигурация бекап сервиса").SetTitleAlign(tview.AlignLeft).SetBorder(true)
}

func updateDirForm(dirForm *tview.Form, pages *tview.Pages, app *tview.Application) {

	dirForm.Clear(true)
	for i, dir := range config.Config.Directories {
		dirForm.AddTextView(
			fmt.Sprintf("Директория #%d", i),
			dir.Dirname,
			40,
			1,
			true,
			false,
		).
			AddButton("Удалить"+dir.Dirname, func() {
				// Логика запуска сервиса
				config.Config.Directories = append(config.Config.Directories[:i], config.Config.Directories[i+1:]...)
				fmt.Println(" Удален...")
			})
	}
	dirForm.AddButton(
		"Главная страница",
		func() {
			// Логика сохранения конфигурации
			pages.SwitchToPage("main")
			fmt.Println(" Настройка создания дампов ДБ.")
		})
	dirForm.AddButton("Добавить", func() {
		config.Config.Directories = append(config.Config.Directories, config.DirectoryConfigType{})
		updateDBForm(dirForm, pages, app)
		app.SetFocus(pages)
	})
	dirForm.SetTitle("Архивация директорий").SetTitleAlign(tview.AlignLeft).SetBorder(true)
}

func updateDBForm(dbForm *tview.Form, pages *tview.Pages, app *tview.Application) {
	dbForm.Clear(true) // Очищаем форму перед обновлением
	for i, db := range config.Config.DataBases {
		dbForm.AddTextView(
			fmt.Sprintf("ДБ #%d", i),
			db.DataBaseName,
			40,
			1,
			true,
			false,
		).
			AddButton("Удалить "+db.DataBaseName, func() {
				config.Config.DataBases = append(config.Config.DataBases[:i], config.Config.DataBases[i+1:]...)
				fmt.Println(" Удален...")
			})
	}
	dbForm.AddButton(
		"Главная страница",
		func() {
			pages.SwitchToPage("main")
			fmt.Println(" Настройка создания дампов ДБ.")
		})
	dbForm.AddButton("Добавить", func() {
		config.Config.DataBases = append(config.Config.DataBases, config.DataBaseConfigType{})
		updateDBForm(dbForm, pages, app)
		app.SetFocus(pages)
	})
	dbForm.SetTitle("Дампы БД").SetTitleAlign(tview.AlignLeft).SetBorder(true)
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
