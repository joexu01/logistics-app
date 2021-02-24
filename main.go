package main

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/joexu01/logistics-app/public"
	"github.com/joexu01/logistics-app/router"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	_ = lib.InitModule("./conf/dev/", []string{"base", "redis"})
	defer lib.Destroy()

	addrs, err := net.InterfaceAddrs()
	if err != nil{
		log.Println(err)
	}
	for _, value := range addrs{
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback(){
			if ipnet.IP.To4() != nil{
				if strings.HasPrefix(ipnet.IP.String(), "192.168") {
					public.LANIPAddr = ipnet.IP.String()
				}
			}
		}
	}

	router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
