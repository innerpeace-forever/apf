package apf

import (
	"fmt"
	"github.com/cihub/seelog"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	config   Configuration
	logger   Logger
	cli      *Cli
	stopChan chan os.Signal
}

type Runner func(*Application) error

func New() *Application {
	config := DefaultConfiguration()
	app := &Application{
		config:   config,
		stopChan: make(chan os.Signal),
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

func (app *Application) WaitStopSignal() {
	signal.Notify(app.stopChan, syscall.SIGINT, syscall.SIGTERM)
	s := <-app.stopChan
	seelog.Infof("get stop signal[%v]", s)
}

func (app *Application) NotifyStopSignal() {
	app.stopChan <- syscall.SIGTERM
	seelog.Info("notify stop signal")
}
