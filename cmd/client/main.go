package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lomoval/otus-golang-project-sysmon/internal/client"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var host string
var port string
var avgInterval string
var notifyInterval string
var metricName string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "connection host")
	flag.StringVar(&port, "port", "8006", "connection port")
	flag.StringVar(&avgInterval, "avgInterval", "5", "interval to calc average (sec)")
	flag.StringVar(&notifyInterval, "notifyInterval", "5", "interval to get notifications (sec)")
	flag.StringVar(&metricName, "metricName", "cpu", "metric name to collect")
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	flag.Parse()

	avgInt, err := strconv.Atoi(avgInterval)
	if err != nil {
		fmt.Printf("incorrect parameter, intervals must be a number: %s", err)
	}
	notifyInt, err := strconv.Atoi(notifyInterval)
	if err != nil {
		fmt.Printf("incorrect parameter, intervals must be a number: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c := client.New(host, port)
	if err := c.Start(ctx, metricName, avgInt, notifyInt); err != nil {
		fmt.Printf("client failed: %s", err)
	}
}
