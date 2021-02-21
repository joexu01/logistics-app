package main

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/joexu01/logistics-app/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = lib.InitModule("./conf/dev/", []string{"base", "redis"})
	defer lib.Destroy()
	router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
