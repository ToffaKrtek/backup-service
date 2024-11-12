package tui

import (
	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var selectedDataBase string
var editDataBase = false

func DatabaseFrame() *tview.Frame {
	config.LoadConfig()
	if len(config.Config.DataBases) < 1 {
		dbForm := tview.NewForm()
		dbForm.AddButton("Создать", func() {
			config.CreateDatabase()
			SetRightItem(2)
		})

		return tview.NewFrame(dbForm)
	} else {
		if editDataBase {
			if db, exist := config.Config.DataBases[selectedDataBase]; exist {
				editDbForm := tview.NewForm()
				editDbForm.
					AddInputField("Пользователь", db.User, 20, nil, func(str string) {
						db.User = str
					}).
					AddPasswordField("Пароль", db.Password, 20, '*', func(str string) {
						db.Password = str
					}).
					AddInputField("Хост", db.Address, 40, nil, func(str string) {
						db.Address = str
					}).
					AddInputField("Имя контейнера", db.ContainerName, 20, nil, func(str string) {
						db.ContainerName = str
					}).
					AddInputField("Имя БД", db.DataBaseName, 40, nil, func(str string) {
						db.DataBaseName = str
					}).
					AddCheckbox("Докер", db.IsDocker, func(checked bool) {
						db.IsDocker = checked
					}).
					AddInputField("Бакет (S3)", db.Bucket, 40, nil, func(str string) {
						db.Bucket = str
					}).
					AddDropDown(
						"Тип БД",
						config.DbTypes,
						config.DbTypesMap[db.TypeDB],
						func(option string, optionIndex int) {
							if optionIndex > 0 {
								db.TypeDB = option
							}
						},
					).
					AddButton("Сохранить", func() {
						editDataBase = false
						config.UpdateDatabase(db)
						SetRightItem(2)
					}).
					AddButton("Удалить", func() {
						editDataBase = false
						config.DeleteDatabase(db)
						SetRightItem(2)
					}).
					AddButton("Отмена", func() {
						editDataBase = false
						SetRightItem(2)
					})
				return tview.NewFrame(editDbForm)
			}
		}
	}
	dbTable := tview.NewTable().SetBorders(true)
	headers := []string{"Наименование", "Тип", "Индекс"}
	for ih, header := range headers {
		dbTable.SetCell(0, ih, tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft))
	}
	i := 0
	for _, db := range config.Config.DataBases {
		dbTable.SetCell(i+1, 0,
			tview.NewTableCell(db.DataBaseName)).
			SetCell(i+1, 1,
				tview.NewTableCell(db.TypeDB)).
			SetCell(i+1, 2,
				tview.NewTableCell(db.Index)) //TODO:: добавить
		i++
	}
	dbTable.SetCell(len(config.Config.DataBases)+1, 0,
		tview.NewTableCell("Новая БД"))

	dbTable.Select(1, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			dbTable.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row, column int) {
		if row == len(config.Config.DataBases)+1 {
			config.CreateDatabase()
			SetRightItem(2)
		} else {
			cell := dbTable.GetCell(row, 2)
			selectedDataBase = cell.Text
			editDataBase = true
			SetRightItem(2)
		}
	})

	return tview.NewFrame(dbTable)
}
