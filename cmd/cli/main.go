package main

import (
	"flag"

	"github.com/ToffaKrtek/backup-service/internal/config"
	"github.com/ToffaKrtek/backup-service/internal/tui"
)

func main() {
	isTmp := flag.Bool("tmp", false, "запустить конфигурацию и историю в /tmp")
	flag.Parse()
	if *isTmp {
		config.SetTmp(*isTmp)
	}

	config.LoadConfig()

	tui.Run()
}
