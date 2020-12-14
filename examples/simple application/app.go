package main

import (
	"github.com/cihub/seelog"
	"github.com/innerpeace-forever/apf"
)

func main() {
	app := apf.New().Configure(
		apf.WithConfiguration(apf.TOML("./examples/simple application/conf.toml")),
		apf.WithLogger(),
		apf.WithCli(apf.NewCli("Test")))

	err := app.Run(func(app *apf.Application) error {
		seelog.Info("Running")
		return nil
	})

	if err != nil {
		seelog.Info("Run Failed! %v", err)
	}
}
