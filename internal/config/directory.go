package config

import "fmt"

type DirectoryConfigType struct {
	Path       string `json:"path"`
	Dirname    string `json:"dirname"`
	Bucket     string `json:"s3_bucket"`
	ObjectName string `json:"object_name"`
	Index      string `json:"index"`
	IsFull     bool   `json:"is_full"`
	Days       int    `json:"days"`
}

func (d DirectoryConfigType) String() string {
	return fmt.Sprintf("DirectoryConfigType{Path: %s, Bucket: %s}", d.Path, d.Bucket)
}

func CreateDirectory() {
	LoadConfig()
	index := getIndex()
	Config.Directories[index] = DirectoryConfigType{
		Index: index,
	}
	SaveConfig(false)
}

func UpdateDirectory(dir DirectoryConfigType) {
	LoadConfig()
	for i, schedule := range Config.Schedules {
		if _, exist := schedule.Directories[dir.Index]; exist {
			Config.Schedules[i].Directories[dir.Index] = dir
		}
	}
	if _, exist := Config.Directories[dir.Index]; exist {
		Config.Directories[dir.Index] = dir
	}
	SaveConfig(false)
}

func DeleteDirectory(dir DirectoryConfigType) {
	LoadConfig()
	for i := range Config.Schedules {
		delete(Config.Schedules[i].Directories, dir.Index)
	}
	delete(Config.Directories, dir.Index)
	SaveConfig(false)
}
