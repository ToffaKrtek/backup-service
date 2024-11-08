package database

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/ToffaKrtek/backup-service/internal/config"
)

type CmdInterface interface {
	Run() error
	SetEnv(env []string)
	SetStdout(stdout io.Writer)
}

type CmdInner struct {
	cmd *exec.Cmd
}

func (e *CmdInner) Run() error {
	return e.cmd.Run()
}

func (e *CmdInner) SetEnv(env []string) {
	e.cmd.Env = env
}

func (e *CmdInner) SetStdout(stdout io.Writer) {
	e.cmd.Stdout = stdout
}

func NewCommand(name string, arg ...string) CmdInterface {
	if execCommand == nil {
		return &CmdInner{
			cmd: exec.Command(name, arg...),
		}
	}
	return execCommand
}

var execCommand CmdInterface

func Dump(wg *sync.WaitGroup, files chan config.S3Item) {
	conf := config.Config
	fmt.Printf("Запуск создания дампов для %d БД", len(conf.DataBases))
	for _, db := range conf.DataBases {
		wg.Add(1)
		go func(db config.DataBaseConfigType) {
			defer wg.Done()
			defer fmt.Println("Закончена архивация")
			switch db.TypeDB {
			case "postgre":
				files <- config.S3Item{
					ObjectName: db.DataBaseName,
					Bucket:     db.Bucket,
					FilePath: dumpPostgreSQLDocker(
						db.ContainerName,
						db.DataBaseName,
						db.User,
						db.Password,
					),
				}
			case "mysql":
				if db.IsDocker {
					item := config.S3Item{
						ObjectName: db.DataBaseName,
						Bucket:     db.Bucket,
						FilePath: dumpMysqlDocker(
							db.ContainerName,
							db.DataBaseName,
							db.User,
							db.Password,
						),
					}
					fmt.Println(item.FilePath)
					files <- item
				} else {
					files <- config.S3Item{
						ObjectName: db.DataBaseName,
						Bucket:     db.Bucket,
						FilePath: dumpMysqlHost(
							db.DataBaseName,
							db.User,
							db.Password,
						),
					}
				}
			}
		}(db)
	}
}

func dumpMysqlDocker(
	containerName string,
	dbName string,
	user string,
	pass string,
) string {
	cmd := NewCommand(
		"docker",
		"exec",
		containerName,
		"mysqldump",
		"-u"+user,
		"-p"+pass,
		dbName,
	)
	filename := config.GetFileName("mysql_dump_%s.sql")
	outfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return ""
	}
	defer outfile.Close()
	cmd.SetStdout(outfile)
	if err := cmd.Run(); err != nil {
		fmt.Println("Ошибка выполнения команды:", err)
		return ""
	}

	return filename

}
func dumpMysqlHost(dbName string, user string, pass string) string {
	cmd := NewCommand(
		"mysqldump",
		"-u"+user,
		"-p"+pass,
		dbName,
	)

	filename := config.GetFileName("mysql_dump_%s.sql")
	outfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return ""
	}
	defer outfile.Close()

	cmd.SetStdout(outfile)
	if err := cmd.Run(); err != nil {
		fmt.Println("Ошибка выполнения команды:", err)
		return ""
	}
	return filename
}

func dumpPostgreSQLDocker(
	containerName string,
	dbName string,
	user string,
	pass string,
) string {
	cmd := NewCommand(
		"docker",
		"exec",
		containerName,
		"pg_dump",
		"-U"+user,
		dbName,
	)

	filename := config.GetFileName("postgresql_dump_%s.sql")
	outfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return ""
	}
	defer outfile.Close()

	cmd.SetStdout(outfile)
	cmd.SetEnv(append(os.Environ(), "PGPASSWORD="+pass))
	if err := cmd.Run(); err != nil {
		fmt.Println("Ошибка выполнения команды:", err)
		return ""
	}

	return filename
}
