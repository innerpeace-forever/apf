package apf

import (
	"fmt"
	"github.com/cihub/seelog"
)

type Application struct {
	config Configuration
	logger Logger
	cli    *Cli
}

type Runner func(*Application) error

func New() *Application {
	config := DefaultConfiguration()
	app := &Application{
		config: config,
	}
	return app
}

func (app *Application) Configure(cfs ...Configurator) *Application {
	for _, cfg := range cfs {
		if cfg != nil {
			cfg(app)
		}
	}
	return app
}

func (app *Application) Run(runner Runner, cfs ...Configurator) error {
	if app.cli != nil {
		if err := app.cli.rootCmd.Execute(); err != nil {
			panic(fmt.Errorf("WithCli Failed! %v", err))
		}
	}

	app.Configure(cfs...)
	err := runner(app)
	if err != nil {
		seelog.Info("App Run Failed! %v", err)
	}

	app.Flush()
	return err
}

func (app *Application) Flush() {
	for _, logger := range app.logger.loggers {
		logger.Flush()
	}
}

func (app *Application) Logger(str string) seelog.LoggerInterface {
	return app.logger.loggers[str]
}

func (app *Application) Config() *Configuration {
	return &app.config
}
