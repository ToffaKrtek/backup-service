package config

import "fmt"

var (
	DbTypes    = []string{"Не выбрано", "mysql", "postgre", "mysql-swarm"}
	DbTypesMap = map[string]int{
		"Mysql":       1,
		"PostgreSql":  2,
		"Mysql-Swarm": 3,
		"":            0,
	}
)

type DataBaseConfigType struct {
	User          string `json:"user"`
	Password      string `json:"password"`
	Address       string `json:"address"`
	ContainerName string `json:"container_name"`
	DataBaseName  string `json:"db_name"`
	IsDocker      bool   `json:"is_docker"`
	Bucket        string `json:"s3_bucket"`
	TypeDB        string `json:"type_db"`
	Index         string `json:"index"`
}

func (d DataBaseConfigType) String() string {
	return fmt.Sprintf("DataBaseConfigType{User: %s, Address: %s, ContainerName: %s, IsDocker: %t, Bucket: %s}", d.User, d.Address, d.ContainerName, d.IsDocker, d.Bucket)
}

func CreateDatabase() {
	LoadConfig()
	index := getIndex()
	Config.DataBases[index] = DataBaseConfigType{
		Index: index,
	}
	SaveConfig(false)
}

func UpdateDatabase(db DataBaseConfigType) {
	LoadConfig()
	for i, schedule := range Config.Schedules {
		if _, exist := schedule.DataBases[db.Index]; exist {
			Config.Schedules[i].DataBases[db.Index] = db
		}
	}
	if _, exist := Config.DataBases[db.Index]; exist {
		Config.DataBases[db.Index] = db
	}
	SaveConfig(false)
}

func DeleteDatabase(db DataBaseConfigType) {
	LoadConfig()
	for i := range Config.Schedules {
		delete(Config.Schedules[i].DataBases, db.Index)
	}
	delete(Config.DataBases, db.Index)
	SaveConfig(false)
}
