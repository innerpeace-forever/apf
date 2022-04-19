package apf

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
)

type ILogger interface {
	Flush()
	Error(v ...interface{}) error
	Errorf(format string, params ...interface{}) error
	Info(v ...interface{})
	Infof(format string, params ...interface{})
}

type ICli interface {
	Execute() error
}

type Application struct {
	conf          *Configuration
	loggers       map[string]ILogger
	loggerCurrent ILogger
	cli           ICli
	stopChan      chan os.Signal
}

type Runner func(*Application) error

func (p *Application) Run(runner Runner) error {
	if p.cli != nil && !reflect.ValueOf(p.cli).IsNil() {
		if err := p.cli.Execute(); err != nil {
			panic(p.Errorf("WithCli Failed! %v", err))
		}
	}

	err := runner(p)
	if err != nil {
		p.Infof("App Run Failed! %v", err)
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
	if p.loggers != nil {
		return p.loggers[str]
	} else {
		return nil
	}
}

func (p *Application) Conf() *Configuration {
	return p.conf
}

func (p *Application) WaitStopSignal() {
	signal.Notify(p.stopChan, syscall.SIGINT, syscall.SIGTERM)
	s := <-p.stopChan
	p.Infof("get stop signal[%v]", s)
}

func (p *Application) NotifyStopSignal() {
	p.stopChan <- syscall.SIGTERM
	p.Info("notify stop signal")
}

// ---------------------------------------------------------------------------
// -------------------------For configuring app-------------------------------

func (p *Application) WithProcFactor(factor int) *Application {
	runtime.GOMAXPROCS(factor * runtime.NumCPU())
	return p
}

func (p *Application) WithCli(cli ICli) *Application {
	p.cli = cli
	return p
}

// ---------------------------------------------------------------------------
// -----------------------------For printing----------------------------------

func (p *Application) Info(v ...interface{}) {
	if p.loggerCurrent != nil && !reflect.ValueOf(p.loggerCurrent).IsNil() {
		p.loggerCurrent.Info(v...)
	} else {
		fmt.Print(v...)
	}
}

func (p *Application) Infof(format string, params ...interface{}) {
	if p.loggerCurrent != nil && !reflect.ValueOf(p.loggerCurrent).IsNil() {
		p.loggerCurrent.Infof(format, params...)
	} else {
		fmt.Printf(format, params...)
	}
}

func (p *Application) Error(v ...interface{}) error {
	if !reflect.ValueOf(p.loggerCurrent).IsNil() {
		return p.loggerCurrent.Error(v...)
	} else {
		_, err := fmt.Print(v...)
		return err
	}
}

func (p *Application) Errorf(format string, params ...interface{}) error {
	if !reflect.ValueOf(p.loggerCurrent).IsNil() {
		return p.loggerCurrent.Errorf(format, params...)
	} else {
		return fmt.Errorf(format, params...)
	}
}

// ---------------------------------------------------------------------------
