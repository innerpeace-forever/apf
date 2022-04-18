package apf

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type ILogger interface {
	Flush()
	Error(v ...interface{}) error
	Errorf(format string, params ...interface{}) error
	Info(v ...interface{})
	Infof(format string, params ...interface{})
}

// TODO: wrap CLI into ICli
type ICli interface {
}

type Application struct {
	config        Configuration
	loggers       map[string]ILogger
	loggerCurrent ILogger
	cli           *Cli
	stopChan      chan os.Signal
}

type Runner func(*Application) error

// TODO: Move all With* for Application into Application.go

func (p *Application) Configure(cfs ...Configurator) *Application {
	for _, cfg := range cfs {
		if cfg != nil {
			cfg(p)
		}
	}
	return p
}

func (p *Application) Run(runner Runner) error {
	// TODO: check p.loggerCurrent is nil
	if p.cli != nil {
		if err := p.cli.rootCmd.Execute(); err != nil {
			panic(p.loggerCurrent.Errorf("WithCli Failed! %v", err))
		}
	}

	err := runner(p)
	if err != nil {
		p.loggerCurrent.Info("App Run Failed! %v", err)
	}

	p.Flush()
	return err
}

func (p *Application) Flush() *Application {
	for _, logger := range p.loggers {
		logger.Flush()
	}
	return p
}

func (p *Application) Logger(str string) ILogger {
	return p.loggers[str]
}

func (p *Application) Config() *Configuration {
	return &p.config
}

func (p *Application) WaitStopSignal() {
	signal.Notify(p.stopChan, syscall.SIGINT, syscall.SIGTERM)
	s := <-p.stopChan
	p.loggerCurrent.Infof("get stop signal[%v]", s)
}

func (p *Application) NotifyStopSignal() {
	p.stopChan <- syscall.SIGTERM
	p.loggerCurrent.Info("notify stop signal")
}
