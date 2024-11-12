package config

import (
	"fmt"
	"time"
)

type ScheduleConfigType struct {
	StartTime    time.Time                      `json:"start_time"`
	ScheduleName string                         `json:"schedule_name"`
	Directories  map[string]DirectoryConfigType `json:"directories"`
	DataBases    map[string]DataBaseConfigType  `json:"data_bases"`
	EveryDay     bool                           `json:"every_day"`
	Active       bool                           `json:"active"`
}

func (c ScheduleConfigType) String() string {
	return fmt.Sprintf("ScheduleConfigType{StartTime: %s, ScheduleName: %s, Directories: %v, DataBases: %v, EveryDay: %v}",
		c.StartTime.Format(time.RFC3339), c.ScheduleName, c.Directories, c.DataBases, c.EveryDay)
}

func CreateSchedule() {
	LoadConfig()
	Config.Schedules = append(Config.Schedules, makeSchedule())
	SaveConfig(false)
}

func makeSchedule() ScheduleConfigType {
	startTime := time.Now().Add(24 * time.Hour)
	if Config != nil {
		startTime = Config.StartTime
	}
	return ScheduleConfigType{
		StartTime:    startTime,
		ScheduleName: getIndex(),
		Directories:  map[string]DirectoryConfigType{},
		DataBases:    map[string]DataBaseConfigType{},
		EveryDay:     true,
	}
}

func SetDataBases(i int, indexes []string) {
	LoadConfig()
	newDbs := map[string]DataBaseConfigType{}
	for _, index := range indexes {
		if val, exists := Config.DataBases[index]; exists {
			newDbs[index] = val
		}
	}
	Config.Schedules[i].DataBases = newDbs
	SaveConfig(true)
}

func SetDirectories(i int, indexes []string) {
	LoadConfig()
	newDirs := map[string]DirectoryConfigType{}
	for _, index := range indexes {
		if val, exists := Config.Directories[index]; exists {
			newDirs[index] = val
		}
	}
	Config.Schedules[i].Directories = newDirs
	SaveConfig(true)
}
