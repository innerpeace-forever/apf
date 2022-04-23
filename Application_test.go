package apf

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestApplication_Run(t *testing.T) {
	Convey("Run with nil Cli, Loggers, Conf, stopChan", t, func() {
		app := &Application{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		}
		So(func() {
			_ = app.Run(func(application *Application) error {
				return nil
			})
		}, ShouldNotPanic)
	})
}

func TestApplication_WithConfService(t *testing.T) {
	Convey("Run with WithConfService", t, func() {
		app := GetApplication().WithTOMLConfiguration("testdata/full_config.toml")
		checkFullConfig(app.Conf())
		So(app.service, ShouldBeNil)
		app.WithConfService()
		So(app.service, ShouldNotBeNil)
	})
}
