package apf

import (
	"os"
)

var app *Application = nil

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
	}
	return app
}
