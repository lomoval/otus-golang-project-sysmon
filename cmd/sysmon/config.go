package main

import (
	"fmt"
	"github.com/lomoval/otus-golang-project-sysmon/internal/logger"
	metricloader "github.com/lomoval/otus-golang-project-sysmon/internal/metrics/loader"
	internalgrpc "github.com/lomoval/otus-golang-project-sysmon/internal/server/grpc"
	"strings"

	"github.com/spf13/viper"
)

const envConfigPrefix = "$env:"

type config struct {
	Server            internalgrpc.Config `yaml:"server"`
	Logger            logger.Config
	IgnoreUnavailable bool
	Metrics           metricloader.Config
}

func newConfig(configFile string) (config, error) {
	cfg := config{}
	viper.SetConfigFile(configFile)

	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", "8006")
	viper.SetDefault("logger.level", "WARN")

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, fmt.Errorf("failed to read config %q: %w", configFile, err)
	}
	keys := viper.AllKeys()
	for _, key := range keys {
		env := viper.GetString(key)
		if strings.HasPrefix(env, envConfigPrefix) {
			err := viper.BindEnv(key, env[len(envConfigPrefix):])
			if err != nil {
				return config{}, fmt.Errorf("failed to prepare config: %w", err)
			}
		}
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to decode into config struct: %w", err)
	}
	return cfg, nil
}
