package tui

import (
	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/rivo/tview"
)

func S3Frame() *tview.Frame {
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
			SetRightItem(4)
		}).
		AddButton("Проверить соединение", func() {
			SetRightItem(4)
		})
	return tview.NewFrame(s3Form)
}
