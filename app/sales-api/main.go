package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/ardanlabs/service/app/sales-api/handlers"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/pkg/errors"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

/*
Replace the defaultservermux for debugging
Figure out timeouts
*/

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	// Construct the application logger.
	log, err := logger.New("SALES-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// CPU Quota

	// Make sure the program is using the correct
	// number of threads if a CPU quota is set.
	if _, err := maxprocs.Set(); err != nil {
		log.Errorw("startup", zap.Error(err))
		os.Exit(1)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
		Keys struct {
			API string `conf:"default:need-key,mask"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "copyright information here",
		},
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	expvar.NewString("build").Set(build)
	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Infow("startup", "config", out)

	// =========================================================================
	// Start Debug Service

	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      handlers.APIMux(log),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
