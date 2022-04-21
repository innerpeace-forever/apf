package apf

import (
	"fmt"
	"github.com/cihub/seelog"
)

func (p *Application) WithConfSeelogger() *Application {
	return p.WithMultiSeelogger(p.conf.Logger)
}

func (p *Application) WithConsoleLogger() *Application {
	p.loggerCurrent = seelog.Current
	return p
}

// WithMultiSeelogger : The last logger with Default tag will be the current logger
func (p *Application) WithMultiSeelogger(confs []map[string]interface{}) *Application {
	if confs == nil || len(confs) == 0 {
		p.loggerCurrent = seelog.Current
		return p
	}

	if p.loggers == nil {
		p.loggers = make(map[string]ILogger, len(confs))
	}

	for _, conf := range confs {
		if isNil(conf["Name"]) || isNil(conf["ConfigFile"]) {
			panic(fmt.Errorf("missing required keys"))
		}
		name := conf["Name"].(string)
		seelogConfFile := conf["ConfigFile"].(string)

		isDefault := ""
		if !isNil(conf["Default"]) {
			isDefault = conf["Default"].(string)
		}

		if logger, err := seelog.LoggerFromConfigAsFile(seelogConfFile); err == nil {
			p.loggers[name] = logger
			if isDefault != "" || p.loggerCurrent == nil {
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
