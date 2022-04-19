package apf

import (
	"fmt"
	"github.com/cihub/seelog"
)

func (p *Application) WithDefaultSeelogger() *Application {
	return p.WithMultiSeelogger(p.conf.Other["Logger"].([]map[string]interface{}))
}

// WithMultiSeelogger : The last logger with Default tag will be the current logger
func (p *Application) WithMultiSeelogger(confs []map[string]interface{}) *Application {
	if confs == nil {
		p.loggerCurrent = seelog.Current
		return p
	}

	if p.loggers == nil {
		p.loggers = make(map[string]ILogger, len(confs))
	}

	for _, conf := range confs {
		name := conf["Name"].(string)
		seelogConfFile := conf["ConfigFile"].(string)
		isDefault := conf["Default"].(string)

		if logger, err := seelog.LoggerFromConfigAsFile(seelogConfFile); err == nil {
			p.loggers[name] = logger
			if isDefault != "" {
				p.loggerCurrent = logger
				if err := seelog.ReplaceLogger(logger); err != nil {
					panic(fmt.Errorf("Load Logger[%s] Configure %s ReplaceLogger Failed! Err:%v\n", name, seelogConfFile, err))
				}
			}
		} else {
			panic(fmt.Errorf("Load Logger[%s] Configure %s Failed! Err:%v\n", name, seelogConfFile, err))
		}
	}

	return p
}
