package database

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ToffaKrtek/backup-service/internal/config"
)

func Dump(files chan string) {
	conf := config.Config
	for _, db := range conf.DataBases {
		switch db.TypeDB {
		case "postgre":
			files <- dumpPostgreSQLDocker(
				db.ContainerName,
				db.DataBaseName,
				db.User,
				db.Password,
			)
		case "mysql":
			if db.IsDocker {
				files <- dumpMysqlDocker(
					db.ContainerName,
					db.DataBaseName,
					db.User,
					db.Password,
				)
			} else {
				files <- dumpMysqlHost(
					db.DataBaseName,
					db.User,
					db.Password,
				)
			}
		}
	}
}

func dumpMysqlDocker(
	containerName string,
	dbName string,
	user string,
	pass string,
) string {
	cmd := exec.Command(
		"docker",
		"exec",
		containerName,
		"mysqldump",
		"-u"+user,
		"-p"+pass,
		dbName,
	)
	filename := getFileName("mysql_dump_%s.sql")
	outfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return ""
	}
	defer outfile.Close()
	cmd.Stdout = outfile
	if err := cmd.Run(); err != nil {
		fmt.Println("Ошибка выполнения команды:", err)
		return ""
	}

	return filename

}
func dumpMysqlHost(dbName string, user string, pass string) string {
	cmd := exec.Command(
		"mysqldump",
		"-u"+user,
		"-p"+pass,
		dbName,
	)

	filename := getFileName("mysql_dump_%s.sql")
	outfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return ""
	}
	defer outfile.Close()

	cmd.Stdout = outfile
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
	cmd := exec.Command(
		"docker",
		"exec",
		containerName,
		"pg_dump",
		"-U"+user,
		dbName,
	)

	filename := getFileName("postgresql_dump_%s.sql")
	outfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return ""
	}
	defer outfile.Close()

	cmd.Stdout = outfile
	cmd.Env = append(os.Environ(), "PGPASSWORD="+pass)
	if err := cmd.Run(); err != nil {
		fmt.Println("Ошибка выполнения команды:", err)
		return ""
	}

	return filename
}

func getFileName(template string) string {
	tempDir := os.TempDir()
	timestamp := time.Now().Format("20060102_150405")
	return filepath.Join(tempDir, fmt.Sprintf(template, timestamp))
}
