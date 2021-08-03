package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var build = "develop"

func main() {
	log.Println("starting service", build)

	n := runtime.NumCPU()
	g := runtime.GOMAXPROCS(0)
	log.Println("NumCPU", n, "GOMAX", g)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("stoping service")
}
