package main

import (
	"context"
	"flag"
	"github.com/lomoval/otus-golang-project-sysmon/internal/logger"
	metricloader "github.com/lomoval/otus-golang-project-sysmon/internal/metrics/loader"
	internalgrpc "github.com/lomoval/otus-golang-project-sysmon/internal/server/grpc"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := newConfig(configFile)
	if err != nil {
		log.Errorf("failed to start %v", err)
		return
	}
	err = logger.PrepareLogger(config.Logger)
	if err != nil {
		log.Errorf("failed to start %v", err)
		return
	}
	collectors, err := metricloader.Load(config.Metrics)
	if err != nil {
		log.Errorf("failed to load collectors: %s", err)
		return
	}

	if len(collectors) == 0 {
		log.Warn("no metrics to collect")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	grpcServer := internalgrpc.NewServer(config.Server, collectors)

	log.Info("system monitoring is running...")

	go func() {
		err = grpcServer.Start()
		if err != nil {
			log.Errorf("grpc server failed: %v", err)
			cancel()
			return
		}
	}()

	<-ctx.Done()
	grpcServer.Stop()
}
