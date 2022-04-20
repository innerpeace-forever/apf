package apf

import (
	"context"
	"errors"
	"github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

type FakeServiceForStart struct{}

func TestCreateService(t *testing.T) {
	Convey("CreateService with port 0", t, func() {
		var b = CreateService("nilServer", 0, seelog.Current)
		So(b, ShouldNotBeNil)
	})

	Convey("CreateService panic with nil logger", t, func() {
		So(func() { CreateService("nilServer", 0, nil) }, ShouldPanic)
	})
}

func (p *FakeServiceForStart) ListenAndServe() error {
	return errors.New("ListenAndServe Fake Error")
}

func (p *FakeServiceForStart) Shutdown(_ context.Context) error {
	return nil
}

func TestBaseService_Start(t *testing.T) {
	Convey("Start with ListenAndServe failed", t, func() {
		var b = &Service{
			name: "FaceServiceForStart",
			svr:  &FakeServiceForStart{},
		}
		So(b.isStopped(), ShouldBeTrue)
	})
}

type FakeServiceForStop struct {
	c chan bool
}

func (p *FakeServiceForStop) ListenAndServe() error {
	p.c = make(chan bool)
	<-p.c
	return nil
}

func (p *FakeServiceForStop) Shutdown(_ context.Context) error {
	p.c <- true
	return nil
}

func TestBaseService_Stop(t *testing.T) {
	Convey("Stop service FakeServiceForStop", t, func() {
		var b = &Service{
			name:   "FakeServiceForStop",
			svr:    &FakeServiceForStop{},
			logger: seelog.Current,
		}

		So(b.isStopped(), ShouldBeTrue)
		b.Start()
		So(b.isStopped(), ShouldBeFalse)
		So(b.Stop(), ShouldBeNil)
		So(b.isStopped(), ShouldBeTrue)
	})

	Convey("Stop service real", t, func() {
		var b = CreateService("real service", 12888, seelog.Current)
		So(b.isStopped(), ShouldBeTrue)
		b.Start()
		So(b.isStopped(), ShouldBeFalse)
		So(b.Stop(), ShouldBeNil)
		So(b.isStopped(), ShouldBeTrue)
	})
}

func TestService_HandleFunc(t *testing.T) {
	Convey("HandleFunc after CreateService", t, func() {
		var b = CreateService("real service", 12888, seelog.Current)
		//b.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		//	_, _ = writer.Write([]byte("This is Test"))
		//})
		b.mux.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("This is Test"))
		})
		b.svr = &http.Server{
			Addr:    ":" + strconv.Itoa(12888),
			Handler: h2c.NewHandler(b.mux, &http2.Server{}),
		}
		b.Start()
		resp, err := http.Get("http://127.0.0.1:12888/test")
		So(err, ShouldNotBeNil)
		So(resp, ShouldNotBeNil)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		So(err, ShouldNotBeNil)
		So(string(body), ShouldEqual, "Test is Test")
		So(b.Stop(), ShouldBeNil)
		So(b.isStopped(), ShouldBeTrue)
	})
}
