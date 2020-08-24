package main

import (
	"apf"
	"fmt"
	"github.com/cihub/seelog"
)

func main() {
	fmt.Printf("Start\n")

	app := apf.New().Configure(
		apf.WithConfiguration(apf.TOML("./examples/simple application/conf.toml")),
		apf.WithLogger())

	err := app.Run(func(app *apf.Application) error {
		seelog.Info("Running")
		return nil
	})

	if err != nil {
		seelog.Info("Run Failed! %v", err)
	}
}
