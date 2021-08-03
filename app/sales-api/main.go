package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

func main() {
	log.Println("starting service", build)

	// Make sure the program is using the correct
	// number of threads if a CPU quota is set.
	if _, err := maxprocs.Set(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	n := runtime.NumCPU()
	g := runtime.GOMAXPROCS(0)
	log.Println("NumCPU", n, "GOMAX", g)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("stoping service")
}
