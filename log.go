package main

import (
	"log"
	"os"
)

type ConfigLog struct {
	Env   string
	Path  string
	Level string
	Type  map[string]string
}

type LogGroup struct {
	logs   map[string]*log.Logger
	config ConfigLog
}

func logGroupNew(config *ConfigLog) (*LogGroup, error) {
	logGroup := &LogGroup{
		config: *config,
		logs:   make(map[string]*log.Logger, len(config.Type)),
	}

	for k, v := range config.Type {
		if config.Env == "development" {
			logGroup.logs[k] = log.New(os.Stdout, "[DEV: "+k+"] ", log.LstdFlags)
		} else {
			logfile, err := os.OpenFile(config.Path+v+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return logGroup, err
			}
			logGroup.logs[k] = log.New(logfile, "["+k+"] ", log.LstdFlags)
		}
	}
	return logGroup, nil
}

func (this *LogGroup) Get(name string) *log.Logger {
	return this.logs[name]
}

func (this *LogGroup) Close() error {
	for _, v := range this.logs {
		if err := ((v.Writer()).(*os.File)).Close(); err != nil {
			return err
		}
	}
	return nil
}
