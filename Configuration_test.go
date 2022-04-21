package apf

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func checkFullConfig(c *Configuration) {
	So(c.Service["Port"], ShouldEqual, 1234)
	So(len(c.Logger), ShouldEqual, 2)
	So(c.Logger[0]["Name"], ShouldEqual, "CommonLogger")
	So(c.Logger[0]["ConfigFile"], ShouldEqual, "log_config/common.xml")
	So(c.Logger[1]["Name"], ShouldEqual, "RuntimeLogger")
	So(c.Logger[1]["ConfigFile"], ShouldEqual, "log_config/runtime.xml")
	So(c.Other["TestOther"], ShouldResemble, map[string]interface{}{"Name": "test"})
}

func TestApplication_WithTOMLConfiguration(t *testing.T) {
	Convey("WithTOMLConfiguration by Full Config", t, func() {
		app := GetApplication().WithTOMLConfiguration("testdata/full_config.toml")
		checkFullConfig(app.Conf())
	})
}

func TestApplication_WithAppendConf(t *testing.T) {
	Convey("WithAppendConf by nil Config", t, func() {
		app := GetApplication().WithTOMLConfiguration("testdata/full_config.toml").WithAppendConf(nil)
		checkFullConfig(app.Conf())
		app.WithAppendConf(&Configuration{
			nil,
			nil,
			map[string]interface{}{"NewOther": map[string]interface{}{"Name": "Test1"}},
		})
		checkFullConfig(app.Conf())
		So(app.Conf().Other["NewOther"], ShouldResemble, map[string]interface{}{"Name": "Test1"})
	})
}
