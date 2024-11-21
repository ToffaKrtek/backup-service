package tui

import (
	"strconv"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var selectedDir string
var editDir = false

func DirFrame() *tview.Frame {
	config.LoadConfig()
	if len(config.Config.Directories) < 1 {
		dirForm := tview.NewForm()
		dirForm.AddButton("Создать", func() {
			config.CreateDirectory()
			SetRightItem(3)
		})
		return tview.NewFrame(dirForm)
	} else {
		if editDir {
			if dir, exist := config.Config.Directories[selectedDir]; exist {
				editDirForm := tview.NewForm()
				editDirForm.
					AddInputField("Путь к папке", dir.Path, 40, nil, func(str string) {
						dir.Path = str
					}).
					AddInputField("Наименование", dir.Dirname, 20, nil, func(str string) {
						dir.Dirname = str
					}).
					AddInputField("Бакет (S3)", dir.Bucket, 20, nil, func(str string) {
						dir.Bucket = str
					}).
					AddCheckbox("Полная архивация", dir.IsFull, func(checked bool) {
						dir.IsFull = checked
					})
				if !dir.IsFull {
					editDirForm.AddInputField(
						"Максимальный срок создания|обновления",
						strconv.Itoa(dir.Days),
						20,
						func(txt string, l rune) bool {
							num, err := strconv.ParseInt(txt, 10, 0)
							if err != nil {
								return false
							}
							return num < 30
						},
						func(text string) {
							num, err := strconv.ParseInt(text, 10, 0)
							if err == nil {
								dir.Days = int(num)
							}
						},
					)
				}

				editDirForm.AddButton("Сохранить", func() {
					editDir = false
					config.UpdateDirectory(dir)
					SetRightItem(3)
				}).
					AddButton("Удалить", func() {
						editDir = false
						config.DeleteDirectory(dir)
						SetRightItem(3)
					}).
					AddButton("Отмена", func() {
						editDir = false
						SetRightItem(3)
					})
				return tview.NewFrame(editDirForm)
			}
		}
	}
	dirTable := tview.NewTable().SetBorders(true)
	headers := []string{"Наименование", "Индекс"}
	for ih, header := range headers {
		dirTable.SetCell(0, ih, tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
	}
	i := 0
	for _, dir := range config.Config.Directories {
		dirTable.SetCell(i+1, 0,
			tview.NewTableCell(dir.Dirname)).
			SetCell(i+1, 1,
				tview.NewTableCell(dir.Index)) //TODO:: добавить
		i++
	}
	dirTable.SetCell(len(config.Config.Directories)+1, 0,
		tview.NewTableCell("Новая Директория"))

	dirTable.Select(1, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			dirTable.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row, column int) {
		if row == len(config.Config.Directories)+1 {
			config.CreateDirectory()
			SetRightItem(3)
		} else {
			cell := dirTable.GetCell(row, 1)
			selectedDir = cell.Text
			editDir = true
			SetRightItem(3)
		}
	})

	return tview.NewFrame(dirTable)
}
