package config

import (
	"ca-generator/example_name_of_project/internal/common/types/error_with_codes"
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	path      string
	envReader envReader
}

type envReader interface {
	EnvReadConfig(addr string, cfg interface{}) error
}

func MustLoad(ctx context.Context, configPath string, envReader envReader) *Config {
	operation := "config.MustLoad()"

	cfg := new(Config)
	cfg.envReader = envReader

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		logrus.WithFields(logrus.Fields{
			"config_path": configPath,
		}).WithError(err).Fatal(error_with_codes.ErrorFailedToFindConfig.SetOperation(operation)) // set operation on custom error
	}

	err = envReader.EnvReadConfig(configPath, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"config_path": configPath,
		}).WithError(err).Fatal(error_with_codes.ErrorFailedToReadConfig.SetOperation(operation)) // set operation on custom error
	}

	return cfg
}
