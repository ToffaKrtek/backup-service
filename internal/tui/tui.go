package tui

import (
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
var app *tview.Application

func Run() {
	app = tview.NewApplication().EnableMouse(true)
	flexTop = tview.NewFlex().SetDirection(tview.FlexColumn)

	formMain := tview.NewForm().
		SetHorizontal(true).
		AddTextView("Статус:", "Активно", 0, 0, true, false).
		AddTextView("Следующий запуск:", config.Config.StartTime.Format(time.RFC822), 0, 0, true, false).
		AddButton("Конфигурация запуска", func() {
			SetRightItem(0)
		}).
		AddButton("Джобы", func() {
			SetRightItem(1)
		}).
		AddButton("Базы данных", func() {
			SetRightItem(2)
		}).
		AddButton("Директории", func() {
			SetRightItem(3)
		}).
		AddButton("Подключение к S3", func() {
			SetRightItem(4)
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
	SetRightItem(0)
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

func SetRightItem(i int) {
	flexTop.RemoveItem(frameRight)
	switch i {
	case 1:
		frameRight = ScheduleFrame()
	case 2:
		frameRight = DatabaseFrame()
	case 3:
		frameRight = DirFrame()
	case 4:
		frameRight = S3Frame()
	default:
		frameRight = ConfigureFrame()
	}
	flexTop.AddItem(frameRight, 0, 3, false)
}
