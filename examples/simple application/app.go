package main

import (
	"github.com/cihub/seelog"
	"github.com/innerpeace-forever/apf"
)

func main() {
	app := apf.GetApplication().
		WithConfiguration(apf.TOML("./examples/simple application/conf.toml")).
		WithLogger().
		WithCli(apf.NewCli("Test")).
		WithProcFactor(2)

	err := app.Run(func(app *apf.Application) error {
		seelog.Info("Running")
		return nil
	})

	if err != nil {
		seelog.Info("Run Failed! %v", err)
	}
}
