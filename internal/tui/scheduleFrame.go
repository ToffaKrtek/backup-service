package tui

import (
	"fmt"
	"time"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var selectedSchedule = 0
var editSchedule = false
var editScheduleDb = false
var editScheduleDir = false

func ScheduleFrame() *tview.Frame {
	config.LoadConfig()
	if len(config.Config.Schedules) < 1 {
		dbForm := tview.NewForm()
		dbForm.AddButton("Создать джобу", func() {
			config.CreateSchedule()
			SetRightItem(1) //TODO::
		})

		return tview.NewFrame(dbForm)
	}
	if editSchedule && selectedSchedule < len(config.Config.Schedules) {
		if editScheduleDb {
			return showDatabaseSelection()
		}
		if editScheduleDir {
			return showDirectorySelection()
		}
		return editScheduleForm()
	}
	scheduleTable := tview.NewTable().SetBorders(true)
	headers := []string{"Наименование", "Ежедневно", "Активно"}
	for ih, header := range headers {
		scheduleTable.SetCell(
			0,
			ih,
			tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft),
		)
	}
	for i, sch := range config.Config.Schedules {
		scheduleTable.SetCell(i+1, 0,
			tview.NewTableCell(sch.ScheduleName)).
			SetCell(i+1, 1,
				tview.NewTableCell(config.ActiveMap[sch.EveryDay])).
			SetCell(i+1, 2,
				tview.NewTableCell(config.ActiveMap[sch.Active])) //TODO:: добавить
	}
	scheduleTable.SetCell(len(config.Config.Schedules)+1, 0,
		tview.NewTableCell("Новая джоба"))

	scheduleTable.Select(1, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			scheduleTable.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row, column int) {
		if row == len(config.Config.Schedules)+1 {
			config.CreateSchedule()
			SetRightItem(1)
		} else {
			selectedSchedule = row - 1
			editSchedule = true
			SetRightItem(1)
		}
	})
	return tview.NewFrame(scheduleTable)
}

func editScheduleForm() *tview.Frame {
	editScheduleForm := tview.NewForm().
		AddTextView(
			"Следующий запуск",
			config.Config.Schedules[selectedSchedule].StartTime.Format(time.RFC822),
			0,
			0,
			true,
			false,
		).
		AddInputField(
			"Наименование",
			config.Config.Schedules[selectedSchedule].ScheduleName,
			40,
			nil,
			func(str string) {
				config.Config.Schedules[selectedSchedule].ScheduleName = str
			}).
		AddCheckbox(
			"Активно",
			config.Config.Schedules[selectedSchedule].Active,
			func(checked bool) {
				config.Config.Schedules[selectedSchedule].Active = checked
			},
		).
		AddCheckbox(
			"Ежедневно",
			config.Config.Schedules[selectedSchedule].EveryDay,
			func(checked bool) {
				config.Config.Schedules[selectedSchedule].EveryDay = checked
			},
		).
		AddButton("Выбрать базы данных", func() {
			editScheduleDb = true
			SetRightItem(1)
		}).
		AddButton("Выбрать базы директории", func() {
			editScheduleDir = true
			SetRightItem(1)
		}).
		AddButton("Сохранить", func() {
			editSchedule = false
			config.SaveConfig(true)
			SetRightItem(1)
		}).
		AddButton("Удалить", func() {
			editSchedule = false
			config.Config.Schedules = append(config.Config.Schedules[:selectedSchedule], config.Config.Schedules[selectedSchedule+1:]...)
			config.SaveConfig(true)
			SetRightItem(1)
		}).
		AddButton("Отмена", func() {
			editSchedule = false
			SetRightItem(1)
		})
	return tview.NewFrame(editScheduleForm)
}
func showDatabaseSelection() *tview.Frame {
	selectionForm := tview.NewForm()
	selectedDatabases := make(map[string]bool)
	for index, _ := range config.Config.DataBases {
		val := false
		if _, exist := config.Config.Schedules[selectedSchedule].DataBases[index]; exist {
			val = true
		}
		selectedDatabases[index] = val
	}
	for index, db := range config.Config.DataBases {
		selectionForm.AddCheckbox(
			fmt.Sprintf("БД %s (%s)", db.DataBaseName, db.Index),
			selectedDatabases[index],
			func(checked bool) {
				selectedDatabases[index] = checked
			},
		)
	}
	selectionForm.AddButton("Сохранить", func() {
		indexes := []string{}
		for index, val := range selectedDatabases {
			if val {
				indexes = append(indexes, index)
			}
		}
		config.SetDataBases(selectedSchedule, indexes)
		//TODO:: add edit database var
		editScheduleDb = false
		SetRightItem(1)
	})
	selectionForm.AddButton("Отмена", func() {
		//TODO:: add edit database var
		editScheduleDb = false
		SetRightItem(1)
	})
	return tview.NewFrame(selectionForm)
}

func showDirectorySelection() *tview.Frame {
	selectionForm := tview.NewForm()
	selectedDirs := make(map[string]bool)
	for index, _ := range config.Config.Directories {
		val := false
		if _, exist := config.Config.Schedules[selectedSchedule].Directories[index]; exist {
			val = true
		}
		selectedDirs[index] = val
	}
	for index, dir := range config.Config.Directories {
		selectionForm.AddCheckbox(
			fmt.Sprintf("Директория %s (%s)", dir.Dirname, dir.Index),
			selectedDirs[index],
			func(checked bool) {
				selectedDirs[index] = checked
			},
		)
	}
	selectionForm.AddButton("Сохранить", func() {
		indexes := []string{}
		for index, val := range selectedDirs {
			if val {
				indexes = append(indexes, index)
			}
		}
		config.SetDirectories(selectedSchedule, indexes)
		//TODO:: add edit database var
		editScheduleDir = false
		SetRightItem(1)
	})
	selectionForm.AddButton("Отмена", func() {
		//TODO:: add edit database var
		editScheduleDir = false
		SetRightItem(1)
	})
	return tview.NewFrame(selectionForm)
}
