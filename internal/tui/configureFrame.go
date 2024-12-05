package tui

import (
	"strconv"
	"time"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/rivo/tview"
)

func ConfigureFrame() *tview.Frame {
	mainForm := tview.NewForm()
	mainForm.AddInputField("Имя сервера", config.Config.ServerName, 20, nil, func(text string) {
		config.Config.ServerName = text
	}).
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
			func(text string) {
				num, err := strconv.ParseInt(text, 10, 0)
				if err == nil {
					t := config.Config.StartTime
					now := time.Now()
					t = time.Date(now.Year(), now.Month(), now.Day(), int(num), t.Minute(), t.Second(), 0, t.Location())

					if t.Before(now) {
						t = t.Add(24 * time.Hour)
					}
					config.Config.StartTime = t
				}
			},
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
			func(text string) {
				num, err := strconv.ParseInt(text, 10, 0)
				if err == nil {
					t := config.Config.StartTime
					now := time.Now()
					t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), int(num), t.Second(), 0, t.Location())

					if t.Before(now) {
						t = t.Add(24 * time.Hour)
					}
					config.Config.StartTime = t
				}
			},
		).
		// AddCheckbox("Запуск каждый день (раз в неделю если нет)", config.Config.EveryDay, func(checked bool) {
		// 	config.Config.EveryDay = checked
		// }).
		AddButton("Сохранить", func() {
			// Логика сохранения конфигурации
			config.UpdateConfig(*config.Config, false)
			SetRightItem(0)
		})
	return tview.NewFrame(mainForm)
}
