package apf

type Application struct {
	config Configuration
	logger Logger
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
	app.Configure(cfs...)
	err := runner(app)
	if err != nil {

	}

	app.Flush()
	return err
}

func (app *Application) Flush() {
	for _, logger := range app.logger.loggers {
		logger.Flush()
	}
}
