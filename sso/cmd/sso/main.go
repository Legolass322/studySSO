package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	config "sso/internal/cfg"
	"sso/internal/logger"
	. "sso/pkg/configuration"
	"syscall"
)

func main() {
	var cfg = config.InitConfiguration(Local)

	var log = logger.New(cfg.Env)
	log.Info("Starting app", slog.String("cfg", cfg.String()))

	application := app.New(
		log,
		cfg,
	)

	go application.GRPCApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signalStop := <-stop

	log.Info("Interrupted by", slog.String("signal", signalStop.String()))

	application.GRPCApp.Stop()
	application.Storage.Stop()

	log.Info("Stopped")

	os.Exit(0)
}
