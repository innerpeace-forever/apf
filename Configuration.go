package apf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cihub/seelog"
	"io/ioutil"
	"path/filepath"
)

// Right now just support TOML
type Configuration struct {
	// Defaults to "info". Possible values are:
	// * "disable"
	// * "fatal"
	// * "error"
	// * "warn"
	// * "info"
	// * "debug"
	LogLevel string                 `toml:"LogLevel"`
	Other    map[string]interface{} `toml:"Other"`
}

type Configurator func(app *Application)

func TOML(filename string) Configuration {
	c := DefaultConfiguration()

	tomlAbsPath, err := filepath.Abs(filename)
	if err != nil {
		panic(fmt.Errorf("toml: %w", err))
	}

	data, err := ioutil.ReadFile(tomlAbsPath)
	if err != nil {
		panic(fmt.Errorf("toml: %w", err))
	}

	if _, err := toml.Decode(string(data), &c); err != nil {
		panic(fmt.Errorf("toml: %w", err))
	}

	return c
}

func DefaultConfiguration() Configuration {
	return Configuration{
		LogLevel: "info",
		Other:    make(map[string]interface{}),
	}
}

func WithCli(cli *Cli) Configurator {
	return func(app *Application) {
		app.cli = cli
	}
}

func WithConfiguration(c Configuration) Configurator {
	return func(app *Application) {
		main := app.config
		if v := c.LogLevel; v != "" {
			main.LogLevel = v
		}

		if v := c.Other; len(v) > 0 {
			if main.Other == nil {
				main.Other = make(map[string]interface{}, len(v))
			}
			for key, value := range v {
				main.Other[key] = value
			}
		}
	}
}

func WithLogger() Configurator {
	return func(app *Application) {
		other := app.config.Other
		var err error

		for _, logger := range other["Logger"].([]map[string]interface{}) {
			name := logger["Name"].(string)
			loggerCfgFile := logger["ConfigFile"].(string)

			if app.logger.loggers == nil {
				app.logger.loggers = make(map[string]seelog.LoggerInterface)
			}

			if app.logger.loggers[name], err = seelog.LoggerFromConfigAsFile(loggerCfgFile); err != nil {
				panic(fmt.Errorf("Load Logger[%s] Configure %s Failed! Err:%v\n", name, loggerCfgFile, err))
			}
		}

		if err = seelog.ReplaceLogger(app.logger.loggers["RuntimeLogger"]); err != nil {
			panic(fmt.Errorf("ReplaceLogger RuntimeLogger Failed! %v", err))
		}
	}
}
