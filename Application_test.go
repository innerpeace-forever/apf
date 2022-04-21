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
