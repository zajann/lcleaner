package config

import (
	"fmt"

	log "github.com/zajann/lcleaner/pkg/easylog"

	"github.com/jinzhu/configor"
)

type Config struct {
	Process
	Log
	Targets []Target
}

type Process struct {
	PIDFilePath string `default:"./"`
	PIDFileName string `default:"lcleaner.pid"`
}

type Log struct {
	FilePath string `default:"../log"`
	FileName string `default:"lcleaner.log"`
	Level    int    `default:"0"`
	MaxSize  int    `default:"10"`
}

type Target struct {
	Path   string
	Regexp string
	Period string
}

func Load(filePath string) (*Config, error) {
	config := new(Config)

	if err := configor.Load(config, filePath); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) DumpToLog() {
	log.Info("##############################################################")
	log.Info("#")
	log.Info("# Settings")
	log.Info("#")
	log.Info("##############################################################")
	log.Info("%-25s : %s", "PID File Path", c.PIDFilePath)
	log.Info("%-25s : %s", "PID File Name", c.PIDFileName)
	log.Info("%-25s : %s", "Log File Path", c.Log.FilePath)
	log.Info("%-25s : %s", "Log File Name", c.Log.FileName)
	log.Info("%-25s : %d", "Log Level", c.Log.Level)
	log.Info("%-25s : %d", "Log File Max Size", c.Log.MaxSize)
	for i, t := range c.Targets {
		log.Info("%-25s : %s", fmt.Sprintf("[%d] Target Path", i), t.Path)
		log.Info("%-25s : %s", fmt.Sprintf("[%d] Target Regexp", i), t.Regexp)
		log.Info("%-25s : %s", fmt.Sprintf("[%d] Target Period", i), t.Period)
	}
	log.Info("##############################################################")
}
