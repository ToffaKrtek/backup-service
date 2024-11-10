package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var down = false
var right = false
var frameRight *tview.Frame
var flexTop *tview.Flex

//var formRight *tview.Form

func main() {
	config.LoadConfig()

	app := tview.NewApplication().EnableMouse(true)
	flexTop = tview.NewFlex().SetDirection(tview.FlexColumn)

	formMain := tview.NewForm().
		SetHorizontal(true).
		AddTextView("Статус:", "Активно", 0, 0, true, false).
		AddTextView("Следующий запуск:", config.Config.StartTime.Format(time.RFC822), 0, 0, true, false).
		AddButton("Конфигурация запуска", func() {
			setRightItem(0)
		}).
		AddButton("Базы данных", func() {
			setRightItem(1)
		}).
		AddButton("Директории", func() {
			setRightItem(2)
		}).
		AddButton("Подключение к S3", func() {
			setRightItem(3)
		}).
		AddButton("Выполнить сейчас", func() {
			config.Config.StartTime = time.Now().Add(2 * time.Second)
			config.SaveConfig(true)
		}).
		AddButton("Отключить", func() {

		})

		//flexMain.AddItem(formMain, 0, 1, false).
	frame_main := tview.NewFrame(formMain)

	flexTop.AddItem(frame_main, 0, 1, true)
	flexTop.AddItem(frameRight, 0, 3, false)
	setRightItem(0)
	frame_main.SetBorder(true).SetTitleAlign(tview.AlignLeft).SetTitle(" Основное меню ")

	// -----------
	table := tview.NewTable().SetBorders(true)
	hc := 0
	headers := strings.Split("Дата, Тип, Размер, Доставка", ",")
	for ih := range headers {
		table.SetCell(0, hc, tview.NewTableCell(headers[ih]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
		hc++
	}

	config.LoadHistory()
	for i, upload := range config.History.Uploads {
		if i > 4 {
			break
		}
		table.SetCell(i+1, 0,
			tview.NewTableCell(upload.DateTime.Format(time.RFC822))).
			SetCell(i+1, 1,
				tview.NewTableCell(upload.ItemType)).
			SetCell(i+1, 2,
				tview.NewTableCell(upload.Size)).
			SetCell(i+1, 3,
				tview.NewTableCell(upload.Status))
	}

	table.Select(1, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	})

	// form.AddButton("Save", func() {}).
	// 	AddButton("Cancel", func() {})

	frame_table := tview.NewFrame(table)
	frame_table.SetBorder(true).SetTitleAlign(tview.AlignLeft).SetTitle(" История запусков ")

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(flexTop, 0, 3, true)
	flex.AddItem(frame_table, 0, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlJ:
			down = !down
		case tcell.KeyCtrlK:
			down = !down
		case tcell.KeyCtrlL:
			right = !right
		case tcell.KeyCtrlH:
			right = !right
		}
		switch true {
		case down:
			app.SetFocus(table)
		case right:
			app.SetFocus(frameRight) //back form if not work
		default:
			app.SetFocus(formMain)
		}
		return event
	})
	if err := app.SetRoot(flex, true).SetFocus(formMain).Run(); err != nil {
		panic(err)
	}
}

func setRightItem(i int) {
	flexTop.RemoveItem(frameRight)
	switch i {
	case 1:
		frameRight = databaseFrame()
	case 2:
		frameRight = dirFrame()
	case 3:
		frameRight = s3Frame()
	default:
		frameRight = configureFrame()
	}
	flexTop.AddItem(frameRight, 0, 3, false)
}

func configureFrame() *tview.Frame {
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
		AddCheckbox("Запуск каждый день (раз в неделю если нет)", config.Config.EveryDay, func(checked bool) {
			config.Config.EveryDay = checked
		}).
		AddButton("Сохранить", func() {
			// Логика сохранения конфигурации
			config.SaveConfig(true)
			setRightItem(0)
		})
	return tview.NewFrame(mainForm)
}

var selectedDataBase = 0
var editDataBase = false

func databaseFrame() *tview.Frame {
	config.LoadConfig()
	if len(config.Config.DataBases) < 1 {
		dbForm := tview.NewForm()
		dbForm.AddButton("Создать", func() {
			config.Config.DataBases = append(config.Config.DataBases, config.DataBaseConfigType{})
			config.SaveConfig(true)
			setRightItem(1)
		})

		return tview.NewFrame(dbForm)
	} else {
		if editDataBase && selectedDataBase < len(config.Config.DataBases) {
			editDbForm := tview.NewForm()
			editDbForm.
				AddInputField("Пользователь", config.Config.DataBases[selectedDataBase].User, 20, nil, func(str string) {
					config.Config.DataBases[selectedDataBase].User = str
				}).
				AddPasswordField("Пароль", config.Config.DataBases[selectedDataBase].Password, 20, '*', func(str string) {
					config.Config.DataBases[selectedDataBase].Password = str
				}).
				AddInputField("Хост", config.Config.DataBases[selectedDataBase].Address, 40, nil, func(str string) {
					config.Config.DataBases[selectedDataBase].Address = str
				}).
				AddInputField("Имя контейнера", config.Config.DataBases[selectedDataBase].ContainerName, 20, nil, func(str string) {
					config.Config.DataBases[selectedDataBase].ContainerName = str
				}).
				AddInputField("Имя БД", config.Config.DataBases[selectedDataBase].DataBaseName, 40, nil, func(str string) {
					config.Config.DataBases[selectedDataBase].DataBaseName = str
				}).
				AddCheckbox("Докер", config.Config.DataBases[selectedDataBase].IsDocker, func(checked bool) {
					config.Config.DataBases[selectedDataBase].IsDocker = checked
				}).
				AddInputField("Бакет (S3)", config.Config.DataBases[selectedDataBase].Bucket, 40, nil, func(str string) {
					config.Config.DataBases[selectedDataBase].Bucket = str
				}).
				AddDropDown(
					"Тип БД",
					config.DbTypes,
					config.DbTypesMap[config.Config.DataBases[selectedDataBase].TypeDB],
					func(option string, optionIndex int) {
						if optionIndex > 0 {
							config.Config.DataBases[selectedDataBase].TypeDB = option
						}
					},
				).
				AddCheckbox("Активно", config.Config.DataBases[selectedDataBase].Active, func(checked bool) {
					config.Config.DataBases[selectedDataBase].Active = checked
				}).
				AddButton("Сохранить", func() {
					editDataBase = false
					config.SaveConfig(true)
					setRightItem(1)
				}).
				AddButton("Удалить", func() {
					editDataBase = false
					config.Config.DataBases = append(config.Config.DataBases[:selectedDataBase], config.Config.DataBases[selectedDataBase+1:]...)
					config.SaveConfig(true)
					setRightItem(1)
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
			SetCell(i+1, 2,
				tview.NewTableCell(config.ActiveMap[db.Active])) //TODO:: добавить
	}
	dbTable.SetCell(len(config.Config.DataBases)+1, 0,
		tview.NewTableCell("Новая БД"))

	dbTable.Select(1, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			dbTable.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row, column int) {
		if row == len(config.Config.DataBases)+1 {
			config.Config.DataBases = append(config.Config.DataBases, config.DataBaseConfigType{})
			config.SaveConfig(true)
			setRightItem(1)
		} else {
			selectedDataBase = row - 1
			editDataBase = true
			setRightItem(1)
		}
	})

	return tview.NewFrame(dbTable)
}

var selectedDir = 0
var editDir = false

func dirFrame() *tview.Frame {
	config.LoadConfig()
	if len(config.Config.Directories) < 1 {
		dirForm := tview.NewForm()
		dirForm.AddButton("Создать", func() {
			config.Config.Directories = append(config.Config.Directories, config.DirectoryConfigType{})
			config.SaveConfig(true)
			setRightItem(2)
		})
		return tview.NewFrame(dirForm)
	} else {
		if editDir && selectedDir < len(config.Config.Directories) {
			editDirForm := tview.NewForm()
			editDirForm.
				AddInputField("Путь к папке", config.Config.Directories[selectedDir].Path, 40, nil, func(str string) {
					config.Config.Directories[selectedDir].Path = str
				}).
				AddInputField("Наименование", config.Config.Directories[selectedDir].Dirname, 20, nil, func(str string) {
					config.Config.Directories[selectedDir].Dirname = str
				}).
				AddInputField("Бакет (S3)", config.Config.Directories[selectedDir].Bucket, 20, nil, func(str string) {
					config.Config.Directories[selectedDir].Bucket = str
				}).
				AddCheckbox("Активно", config.Config.Directories[selectedDir].Active, func(checked bool) {
					config.Config.Directories[selectedDir].Active = checked
				}).
				AddButton("Сохранить", func() {
					editDir = false
					config.SaveConfig(true)
					setRightItem(2)
				}).
				AddButton("Удалить", func() {
					editDir = false
					config.Config.Directories = append(config.Config.Directories[:selectedDir], config.Config.Directories[selectedDir+1:]...)
					config.SaveConfig(true)
					setRightItem(2)
				})
			return tview.NewFrame(editDirForm)
		}
	}
	dirTable := tview.NewTable().SetBorders(true)
	headers := strings.Split("Наименование, Статус", ",")
	hc := 0
	for ih := range headers {
		dirTable.SetCell(0, hc, tview.NewTableCell(headers[ih]).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
		hc++
	}
	for i, dir := range config.Config.Directories {
		dirTable.SetCell(i+1, 0,
			tview.NewTableCell(dir.Dirname)).
			SetCell(i+1, 2,
				tview.NewTableCell(config.ActiveMap[dir.Active])) //TODO:: добавить
	}
	dirTable.SetCell(len(config.Config.Directories)+1, 0,
		tview.NewTableCell("Новая Директория"))

	dirTable.Select(1, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			dirTable.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row, column int) {
		if row == len(config.Config.Directories)+1 {
			config.Config.Directories = append(config.Config.Directories, config.DirectoryConfigType{})
			config.SaveConfig(true)
			setRightItem(2)
		} else {
			selectedDir = row - 1
			editDir = true
			setRightItem(2)
		}
	})

	return tview.NewFrame(dirTable)
}

func s3Frame() *tview.Frame {
	s3Form := tview.NewForm().
		AddInputField("URL", config.Config.S3.Endpoint, 40, nil, func(str string) {
			config.Config.S3.Endpoint = str
		}).
		AddInputField("AccessKeyID", config.Config.S3.AccessKeyID, 40, nil, func(str string) {
			config.Config.S3.AccessKeyID = str
		}).
		AddPasswordField("SecretAccessKey", config.Config.S3.SecretAccessKey, 40, '*', func(str string) {
			config.Config.S3.SecretAccessKey = str
		}).
		AddCheckbox("Отправлять бекапы", config.Config.S3.Send, func(checked bool) {
			config.Config.S3.Send = checked
		}).
		AddButton("Сохранить", func() {
			config.SaveConfig(true)
			setRightItem(3)
		}).
		AddButton("Проверить соединение", func() {
			setRightItem(3)
		})
	return tview.NewFrame(s3Form)
}
