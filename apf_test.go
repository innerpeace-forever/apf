package apf

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetApplication(t *testing.T) {
	Convey("GetApplication multi times", t, func() {
		app1 := GetApplication()
		app2 := GetApplication()
		So(app1, ShouldEqual, app2)
	})

	Convey("GetApplication after New", t, func() {
		app1 := New()
		app2 := GetApplication()
		So(app1, ShouldEqual, app2)
	})
}

func TestNew(t *testing.T) {
	Convey("New is nil", t, func() {
		app := New()
		So(app, ShouldNotBeNil)
	})
}
