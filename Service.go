package apf

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"strconv"
	"time"
)

// IServer can be used for fake testing
type IServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type Service struct {
	svr    IServer
	name   string
	port   int
	mux    *http.ServeMux
	logger ILogger
	cStop  chan bool
}

type Context struct {
	Request       *http.Request
	ResponseWrite http.ResponseWriter
}

type Handler func(*Context)

func CreateService(name string, port int, logger ILogger) *Service {
	if isNil(logger) {
		panic(fmt.Errorf("nil input for BaseService Start [%s]", name))
	}

	mux := http.NewServeMux()

	return &Service{
		svr: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: h2c.NewHandler(mux, &http2.Server{}),
		},
		name:   name,
		port:   port,
		mux:    mux,
		logger: logger,
	}
}

func (p *Service) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *Service {
	p.mux.HandleFunc(pattern, handler)
	return p
}

func (p *Service) Handle(pattern string, handler Handler) *Service {
	p.mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		ctx := &Context{
			Request:       request,
			ResponseWrite: writer,
		}

		handler(ctx)
	})
	return p
}

func (p *Service) Start() {
	if p.cStop == nil {
		p.cStop = make(chan bool)
	}

	go func() {
		defer func() {
			p.cStop <- true
		}()

		p.logger.Infof("Server[%s] Started", p.name)
		if err := p.svr.ListenAndServe(); err != nil {
			_ = p.logger.Errorf("Server[%s] ListenAndServe Failed: Err:%v", p.name, err)
		}
	}()
}

func (p *Service) isStopped() bool {
	if p.cStop == nil {
		return true
	}

	select {
	case <-p.cStop:
		return true
	case <-time.After(1 * time.Second):
		return false
	}
}

func (p *Service) Stop() error {
	if p.svr == nil {
		return errors.New("nil Server when Stop")
	}

	p.logger.Infof("Server[%s] Start to Stop", p.name)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := p.svr.Shutdown(ctx); err != nil {
		_ = p.logger.Errorf("Server[%s] Shutdown Error:", p.name, err)
	}

	p.logger.Infof("Server[%s] Shutdown Completed", p.name)
	return nil
}
