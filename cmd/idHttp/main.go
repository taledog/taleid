package main

import (
	"fmt"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"

	klog "github.com/go-kit/kit/log"
	"github.com/taledog/taleid/cmd/idHttp/conf"
	"github.com/taledog/taleid/cmd/idHttp/endpoints"
	"github.com/taledog/taleid/cmd/idHttp/http"
	"github.com/taledog/taleid/cmd/idHttp/service"
	"github.com/taledog/taleid/store"
)

func main() {
	config := conf.Init()
	var logger klog.Logger
	logger = klog.NewLogfmtLogger(klog.NewSyncWriter(os.Stderr))
	logger = klog.With(logger, "ts", klog.DefaultTimestampUTC)

	db := store.NewStore(config.Store.Type, config.Store.URI, config.Generate.DataCenter, logger)
	svc := service.New(logger, db, config.Generate.Step, config.Generate.DataCenter, config.Store.Min, config.Store.Max)
	e := endpoints.New(svc, logger)
	h := http.NewHTTPHandler(e)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", config.Server.Address)
		errs <- netHttp.ListenAndServe(config.Server.Address, h)
	}()

	logger.Log("exit", <-errs)
}
