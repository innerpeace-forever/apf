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
	oldApp := app

	app = &Application{
		conf:     DefaultConfiguration(),
		stopChan: make(chan os.Signal),
		cli:      nil,
	}

	if oldApp != nil {
		oldApp.NotifyStopSignal()
		oldApp.Flush()
	}

	return app
}
