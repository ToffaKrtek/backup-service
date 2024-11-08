package cli

import (
	"strconv"
	"strings"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ConfigureFrame() *tview.Frame {
	mainForm := tview.NewForm()
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
		AddButton("Сохранить", func() {
			// Логика сохранения конфигурации
		})
	return tview.NewFrame(mainForm)
}

var selectedDataBase = 0
var editDataBase = false

func DatabaseFrame() *tview.Frame {
	config.LoadConfig()
	if len(config.Config.DataBases) < 1 {
		dbForm := tview.NewForm()
		dbForm.AddButton("Создать", func() {
			config.Config.DataBases = append(config.Config.DataBases, config.DataBaseConfigType{})
		})

		return tview.NewFrame(dbForm)
	} else {
		if editDataBase && selectedDataBase < len(config.Config.DataBases) {
			editDbForm := tview.NewForm()
			editDbForm.AddInputField("Имя сервера", config.Config.ServerName, 20, nil, nil).
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
				AddButton("Сохранить", func() {
					// Логика сохранения конфигурации
				})
			return tview.NewFrame(editDbForm)
		}
	}
	dbTable := tview.NewTable().SetBorders(true)
	headers := strings.Split("Наименование, Тип, Статус", ",")
	hc := 0
	for ih := range headers {
		dbTable.SetCell(0, hc, tview.NewTableCell(headers[ih]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
		hc++
	}
	for i, db := range config.Config.DataBases {
		dbTable.SetCell(i+1, 0,
			tview.NewTableCell(db.DataBaseName)).
			SetCell(i+1, 1,
				tview.NewTableCell(db.TypeDB)).
			SetCell(i+1, 3,
				tview.NewTableCell("Активно")) //TODO:: добавить
	}

	dbTable.Select(1, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			dbTable.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row, column int) {
		selectedDataBase = row - 1
		editDataBase = true
	})

	return tview.NewFrame(dbTable)
}

func DirFrame() *tview.Frame {
	config.LoadConfig()
	dirForm := tview.NewForm()
	if len(config.Config.DataBases) < 1 {
		//tview.
	}
	return tview.NewFrame(dirForm)
}

func S3Frame() *tview.Frame {
	s3Form := tview.NewForm()
	return tview.NewFrame(s3Form)
}
