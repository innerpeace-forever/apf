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
	global   map[string]interface{}
}

var app *Application = nil

type Runner func(*Application) error

func init() {
	app = GetApplication()
}

func GetApplication() *Application {
	if app != nil {
		return app
	} else {
		return New()
	}
}

func New() *Application {
	config := DefaultConfiguration()
	app := &Application{
		config:   config,
		stopChan: make(chan os.Signal),
		cli:      NewCli("Default CLi"),
		global:   make(map[string]interface{}),
	}
	return app
}

func (p *Application) Configure(cfs ...Configurator) *Application {
	for _, cfg := range cfs {
		if cfg != nil {
			cfg(p)
		}
	}
	return p
}

func (p *Application) Run(runner Runner, cfs ...Configurator) error {
	if p.cli != nil {
		if err := p.cli.rootCmd.Execute(); err != nil {
			panic(fmt.Errorf("WithCli Failed! %v", err))
		}
	}

	p.Configure(cfs...)
	err := runner(p)
	if err != nil {
		seelog.Info("App Run Failed! %v", err)
	}

	p.Flush()
	return err
}

func (p *Application) Flush() *Application {
	for _, logger := range p.logger.loggers {
		logger.Flush()
	}
	return p
}

func (p *Application) Logger(str string) seelog.LoggerInterface {
	return p.logger.loggers[str]
}

func (p *Application) Config() *Configuration {
	return &p.config
}

func (p *Application) GetGlobal(key string) interface{} {
	if v, ok := p.global[key]; ok {
		return v
	} else {
		return nil
	}
}

func (p *Application) SetGlobal(key string, value interface{}) {
	p.global[key] = value
}

func (p *Application) WaitStopSignal() {
	signal.Notify(p.stopChan, syscall.SIGINT, syscall.SIGTERM)
	s := <-p.stopChan
	seelog.Infof("get stop signal[%v]", s)
}

func (p *Application) NotifyStopSignal() {
	p.stopChan <- syscall.SIGTERM
	seelog.Info("notify stop signal")
}
