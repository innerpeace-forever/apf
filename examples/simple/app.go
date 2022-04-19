package main

import (
	"fmt"
	"github.com/innerpeace-forever/apf"
)

func main() {
	fmt.Printf("Start\n")

	if err := apf.GetApplication().
		WithTOMLConfiguration("conf.toml").
		WithDefaultSeelogger().
		Run(func(app *apf.Application) error {
			app.Info("Running")
			return nil
		}); err != nil {
		apf.GetApplication().Infof("Run Failed! %v", err)
	}
}
